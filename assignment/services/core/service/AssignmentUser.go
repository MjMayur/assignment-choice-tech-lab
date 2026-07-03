package service

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"project/entity"
	assignmentContext "project/pkg/context"
	"project/pkg/errors"
	"project/pkg/log"

	"github.com/go-redis/redis/v8"
)

type AssignmentUserServiceImpl struct {
	assignmentUserRepo AssignmentUserRepo
	redisClient        *redis.Client
	cacheTTL           time.Duration
}

func NewAssignmentUserServiceImpl(assignmentUserRepo AssignmentUserRepo, redisClient *redis.Client) (*AssignmentUserServiceImpl, error) {
	return &AssignmentUserServiceImpl{
		assignmentUserRepo: assignmentUserRepo,
		redisClient:        redisClient,
		cacheTTL:           5 * time.Minute,
	}, nil
}

func (s *AssignmentUserServiceImpl) StoreCsvRecords(ctx context.Context, csvFile *multipart.FileHeader) (map[string]interface{}, map[string]interface{}, string, errors.Response) {
	reqID, _ := assignmentContext.GetRequestIDFromContext(ctx)
	log.Info("service: StoreCsvRecords started", reqID)

	file, err := csvFile.Open()
	if err != nil {
		log.Error("Failed to open CSV file: "+err.Error(), reqID)
		return nil, nil, "", errors.ResponseInternalServerError("Failed to open CSV file")
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	header, err := reader.Read()
	if err != nil {
		log.Error("Failed to read CSV header: "+err.Error(), reqID)
		return nil, nil, "", errors.ResponseBadRequestError("Failed to read CSV header: " + err.Error())
	}

	expectedHeaders := []string{"first_name", "last_name", "company_name", "address", "city", "county", "postal", "phone", "email", "web"}
	if err := s.validateCSVHeader(header, expectedHeaders); err != nil {
		log.Error("Invalid CSV header: "+err.Error(), reqID)
		return nil, nil, "", errors.ResponseBadRequestError("Invalid CSV header: " + err.Error())
	}

	cleanRecords := []entity.AssignmentUser{}
	ignoredRecords := []map[string]interface{}{}
	lineNumber := 2

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			ignoredRecords = append(ignoredRecords, map[string]interface{}{
				"line":   lineNumber,
				"error":  "Failed to read row: " + err.Error(),
				"record": record,
			})
			lineNumber++
			continue
		}

		user, isValid, validationErrors := s.validateCSVRecord(record, header)
		if isValid {
			cleanRecords = append(cleanRecords, user)
		} else {
			ignoredRecords = append(ignoredRecords, map[string]interface{}{
				"line":   lineNumber,
				"errors": validationErrors,
				"record": record,
			})
		}
		lineNumber++
	}

	var insertedCount int
	var successMsg string
	if len(cleanRecords) > 0 {
		insertedCount, err = s.assignmentUserRepo.BulkInsertAssignmentUsers(ctx, cleanRecords)
		if err != nil {
			log.Error("Failed to insert records: "+err.Error(), reqID)
			return nil, nil, "", errors.ResponseInternalServerError("Failed to insert records: " + err.Error())
		}
		s.refreshAssignmentUserListCache(ctx)
		successMsg = fmt.Sprintf("Successfully inserted %d records", insertedCount)
	} else {
		successMsg = "No valid records to insert"
	}

	storeMap := map[string]interface{}{
		"totalRecords":  len(cleanRecords),
		"insertedCount": insertedCount,
		"records":       cleanRecords,
	}

	ignoreMap := map[string]interface{}{
		"totalIgnored": len(ignoredRecords),
		"records":      ignoredRecords,
	}

	log.Info(fmt.Sprintf("StoreCsvRecords completed - Inserted: %d, Ignored: %d", insertedCount, len(ignoredRecords)), reqID)
	return storeMap, ignoreMap, successMsg, nil
}

func (s *AssignmentUserServiceImpl) validateCSVHeader(header []string, expectedHeaders []string) error {
	if len(header) < len(expectedHeaders) {
		return fmt.Errorf("expected at least %d columns, got %d", len(expectedHeaders), len(header))
	}

	headerMap := make(map[string]bool)
	for _, h := range header {
		headerMap[strings.ToLower(strings.TrimSpace(h))] = true
	}

	var missingColumns []string
	for _, expected := range expectedHeaders {
		if !headerMap[expected] {
			missingColumns = append(missingColumns, expected)
		}
	}

	if len(missingColumns) > 0 {
		return fmt.Errorf("missing required columns: %v", missingColumns)
	}

	return nil
}

func (s *AssignmentUserServiceImpl) validateCSVRecord(record []string, header []string) (entity.AssignmentUser, bool, []string) {
	var validationErrors []string
	user := entity.AssignmentUser{}

	headerMap := make(map[string]int)
	for i, h := range header {
		headerMap[strings.ToLower(strings.TrimSpace(h))] = i
	}

	getValue := func(columnName string) string {
		if idx, ok := headerMap[columnName]; ok && idx < len(record) {
			return strings.TrimSpace(record[idx])
		}
		return ""
	}

	user.FirstName = getValue("first_name")
	user.LastName = getValue("last_name")
	user.CompanyName = getValue("company_name")
	user.Address = getValue("address")
	user.City = getValue("city")
	user.County = getValue("county")
	user.Postal = getValue("postal")
	user.Phone = getValue("phone")
	user.Email = getValue("email")
	user.Web = getValue("web")

	if user.FirstName == "" {
		validationErrors = append(validationErrors, "first_name is required")
	}
	if user.LastName == "" {
		validationErrors = append(validationErrors, "last_name is required")
	}
	if user.Email != "" && !s.isValidEmailFormat(user.Email) {
		validationErrors = append(validationErrors, "invalid email format: "+user.Email)
	}
	if user.Phone != "" && len(user.Phone) < 5 {
		validationErrors = append(validationErrors, "phone number is too short (minimum 5 characters)")
	}
	if user.Postal != "" && len(user.Postal) < 2 {
		validationErrors = append(validationErrors, "postal code is too short (minimum 2 characters)")
	}

	return user, len(validationErrors) == 0, validationErrors
}

func (s *AssignmentUserServiceImpl) isValidEmailFormat(email string) bool {
	if email == "" {
		return true
	}
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	localPart := parts[0]
	domainPart := parts[1]
	if localPart == "" || domainPart == "" || !strings.Contains(domainPart, ".") {
		return false
	}

	domainParts := strings.Split(domainPart, ".")
	if len(domainParts) < 2 || domainParts[len(domainParts)-1] == "" || len(domainParts[len(domainParts)-1]) < 2 {
		return false
	}

	return true
}

func (s *AssignmentUserServiceImpl) CreateAssignmentUser(ctx context.Context, assignmentUser entity.AssignmentUser) (int, errors.Response) {
	reqID, _ := assignmentContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>assignmentUser: create assignment user started", reqID)

	assignmentUserID, errResp := s.assignmentUserRepo.CreateAssignmentUser(ctx, assignmentUser)
	if errResp != nil {
		return 0, errResp
	}

	s.refreshAssignmentUserListCache(ctx)
	log.Info("core>service>assignmentUser: create assignment user completed", reqID)
	return assignmentUserID, nil
}

func (s *AssignmentUserServiceImpl) PartialUpdateAssignmentUser(ctx context.Context, assignmentUser entity.AssignmentUser) errors.Response {
	reqID, _ := assignmentContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>assignmentUser: partial update assignment user started", reqID)

	errResp := s.assignmentUserRepo.PartialUpdateAssignmentUser(ctx, assignmentUser)
	if errResp != nil {
		return errResp
	}

	s.refreshAssignmentUserListCache(ctx)
	log.Info("core>service>assignmentUser: partial update assignment user completed", reqID)
	return nil
}

func (s *AssignmentUserServiceImpl) DeleteAssignmentUser(ctx context.Context, assignmentUserID int) errors.Response {
	reqID, _ := assignmentContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>assignmentUser: delete assignment user started", reqID)

	errResp := s.assignmentUserRepo.DeleteAssignmentUser(ctx, assignmentUserID)
	if errResp != nil {
		return errResp
	}

	s.refreshAssignmentUserListCache(ctx)
	log.Info("core>service>assignmentUser: delete assignment user completed", reqID)
	return nil
}

func (s *AssignmentUserServiceImpl) GetAssignmentUserbyID(ctx context.Context, assignmentUserID int) (entity.AssignmentUser, errors.Response) {
	reqID, _ := assignmentContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>assignmentUser: get assignment user started", reqID)

	cacheKey := fmt.Sprintf("assignment_users:detail:%d", assignmentUserID)
	var cachedUser entity.AssignmentUser
	if s.getCachedValue(ctx, cacheKey, &cachedUser) {
		return cachedUser, nil
	}

	assignmentUser, errResp := s.assignmentUserRepo.GetAssignmentUserbyID(ctx, assignmentUserID)
	if errResp != nil {
		return assignmentUser, errResp
	}

	s.setCachedValue(ctx, cacheKey, assignmentUser)
	log.Info("core>service>assignmentUser: get assignment user completed", reqID)
	return assignmentUser, nil
}

func (s *AssignmentUserServiceImpl) ListAssignmentUser(ctx context.Context) (int, []entity.AssignmentUser, errors.Response) {
	reqID, _ := assignmentContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>assignmentUser: assignment user list started", reqID)

	cacheKey := "assignment_users:list"
	var cachedList struct {
		TotalRecords int                     `json:"totalRecords"`
		Records      []entity.AssignmentUser `json:"records"`
	}
	if s.getCachedValue(ctx, cacheKey, &cachedList) {
		return cachedList.TotalRecords, cachedList.Records, nil
	}

	totalRecords, assignmentUsers, errResp := s.assignmentUserRepo.ListAssignmentUser(ctx)
	if errResp != nil {
		return 0, nil, errResp
	}

	cachedList = struct {
		TotalRecords int                     `json:"totalRecords"`
		Records      []entity.AssignmentUser `json:"records"`
	}{
		TotalRecords: totalRecords,
		Records:      assignmentUsers,
	}
	s.setCachedValue(ctx, cacheKey, cachedList)
	log.Info("core>service>assignmentUser: assignment user list completed", reqID)
	return totalRecords, assignmentUsers, nil
}

func (s *AssignmentUserServiceImpl) refreshAssignmentUserListCache(ctx context.Context) {
	if s.redisClient == nil {
		return
	}

	totalRecords, assignmentUsers, errResp := s.assignmentUserRepo.ListAssignmentUser(ctx)
	if errResp != nil {
		log.Error("failed to refresh assignment user list cache: "+errResp.Error(), "")
		return
	}

	cachedList := struct {
		TotalRecords int                     `json:"totalRecords"`
		Records      []entity.AssignmentUser `json:"records"`
	}{
		TotalRecords: totalRecords,
		Records:      assignmentUsers,
	}

	if err := s.setCachedValue(ctx, "assignment_users:list", cachedList); err != nil {
		log.Error("failed to set assignment user list cache: "+err.Error(), "")
	}
}

func (s *AssignmentUserServiceImpl) getCachedValue(ctx context.Context, key string, target interface{}) bool {
	if s.redisClient == nil {
		return false
	}

	payload, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		return false
	}

	if err := json.Unmarshal([]byte(payload), target); err != nil {
		log.Error("failed to decode cache payload: "+err.Error(), "")
		return false
	}

	return true
}

func (s *AssignmentUserServiceImpl) setCachedValue(ctx context.Context, key string, value interface{}) error {
	if s.redisClient == nil {
		return nil
	}

	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return s.redisClient.Set(ctx, key, payload, s.cacheTTL).Err()
}
