package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"project/entity"

	assignmentContext "project/pkg/context"
	"project/pkg/errors"
	"project/pkg/log"
	"strconv"
	"strings"
	"time"

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

// StoreCsvRecords processes CSV file and stores records
// StoreCsvRecords processes CSV file and stores records
func (s *AssignmentUserServiceImpl) StoreCsvRecords(ctx context.Context, csvFile *multipart.FileHeader) (map[string]interface{}, map[string]interface{}, string, errors.Response) {
	reqID, _ := assignmentContext.GetRequestIDFromContext(ctx)
	log.Info("service: StoreCsvRecords started", reqID)

	// Open the uploaded file
	file, err := csvFile.Open()
	if err != nil {
		log.Error("Failed to open CSV file: "+err.Error(), reqID)
		return nil, nil, "", errors.ResponseInternalServerError("Failed to open CSV file")
	}
	defer file.Close()

	// Create CSV reader
	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	// Read header
	header, err := reader.Read()
	if err != nil {
		log.Error("Failed to read CSV header: "+err.Error(), reqID)
		return nil, nil, "", errors.ResponseBadRequestError("Failed to read CSV header: " + err.Error())
	}

	// Validate header
	expectedHeaders := []string{"first_name", "last_name", "company_name", "address", "city", "county", "postal", "phone", "email", "web"}
	if err := s.validateCSVHeader(header, expectedHeaders); err != nil {
		log.Error("Invalid CSV header: "+err.Error(), reqID)
		return nil, nil, "", errors.ResponseBadRequestError("Invalid CSV header: " + err.Error())
	}

	// Process records
	cleanRecords := []entity.AssignmentUser{}
	ignoredRecords := []map[string]interface{}{}
	lineNumber := 2 // Start from line 2 (after header)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Info(fmt.Sprintf("Error reading CSV row %d: %v", lineNumber, err), reqID)
			ignoredRecords = append(ignoredRecords, map[string]interface{}{
				"line":   lineNumber,
				"error":  "Failed to read row: " + err.Error(),
				"record": record,
			})
			lineNumber++
			continue
		}

		// Validate record
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

	log.Info(fmt.Sprintf(fmt.Sprintf("Validation completed - Clean: %d, Ignored: %d", len(cleanRecords), len(ignoredRecords))), reqID)

	// Insert clean records
	var insertedCount int
	var successMsg string
	if len(cleanRecords) > 0 {
		insertedCount, err = s.assignmentUserRepo.BulkInsertAssignmentUsers(ctx, cleanRecords)
		if err != nil {
			log.Error("Failed to insert records: "+err.Error(), reqID)
			return nil, nil, "", errors.ResponseInternalServerError("Failed to insert records: " + err.Error())
		}
		successMsg = fmt.Sprintf("Successfully inserted %d records", insertedCount)
	} else {
		successMsg = "No valid records to insert"
	}

	// Prepare response maps
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

// validateCSVHeader validates CSV header
func (s *AssignmentUserServiceImpl) validateCSVHeader(header []string, expectedHeaders []string) error {
	if len(header) < len(expectedHeaders) {
		return fmt.Errorf("expected at least %d columns, got %d", len(expectedHeaders), len(header))
	}

	// Convert header to lowercase
	headerMap := make(map[string]bool)
	for _, h := range header {
		headerMap[strings.ToLower(strings.TrimSpace(h))] = true
	}

	// Check required columns
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

// validateCSVRecord validates a single CSV record
func (s *AssignmentUserServiceImpl) validateCSVRecord(record []string, header []string) (entity.AssignmentUser, bool, []string) {
	var errors []string
	user := entity.AssignmentUser{}

	// Create header map
	headerMap := make(map[string]int)
	for i, h := range header {
		headerMap[strings.ToLower(strings.TrimSpace(h))] = i
	}

	// Helper to get value
	getValue := func(columnName string) string {
		if idx, ok := headerMap[columnName]; ok && idx < len(record) {
			return strings.TrimSpace(record[idx])
		}
		return ""
	}

	// Extract values
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

	// Validate First Name
	if user.FirstName == "" {
		errors = append(errors, "first_name is required")
	}

	// Validate Last Name
	if user.LastName == "" {
		errors = append(errors, "last_name is required")
	}

	// Validate Email (if provided)
	if user.Email != "" && !s.isValidEmailFormat(user.Email) {
		errors = append(errors, "invalid email format: "+user.Email)
	}

	// Validate Phone (if provided)
	if user.Phone != "" && len(user.Phone) < 5 {
		errors = append(errors, "phone number is too short (minimum 5 characters)")
	}

	// Validate Postal code (if provided)
	if user.Postal != "" && len(user.Postal) < 2 {
		errors = append(errors, "postal code is too short (minimum 2 characters)")
	}

	isValid := len(errors) == 0
	return user, isValid, errors
}

// isValidEmailFormat validates email format
func (s *AssignmentUserServiceImpl) isValidEmailFormat(email string) bool {
	if email == "" {
		return true
	}

	// Check if contains @ and .
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}

	// Split email into local and domain parts
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	localPart := parts[0]
	domainPart := parts[1]

	// Check if local part is not empty
	if localPart == "" {
		return false
	}

	// Check if domain part has at least one dot and is not empty
	if domainPart == "" || !strings.Contains(domainPart, ".") {
		return false
	}

	// Check if domain part has valid TLD (at least 2 characters after last dot)
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
		log.Error(errResp.Error(), reqID)
		return 0, errResp
	}

	log.Info("core>service>assignmentUser: create assignment user completed & assignment user id is "+strconv.Itoa(assignmentUserID), reqID)
	return assignmentUserID, nil
}

func (s *AssignmentUserServiceImpl) PartialUpdateAssignmentUser(ctx context.Context, assignmentUser entity.AssignmentUser) errors.Response {
	reqID, _ := assignmentContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>assignmentUser: partila update assignment user started for assignment user id "+strconv.Itoa(assignmentUser.ID), reqID)

	errResp := s.assignmentUserRepo.PartialUpdateAssignmentUser(ctx, assignmentUser)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return errResp
	}

	log.Info("core>service>assignmentUser: update assignment user completed for assignment user id "+strconv.Itoa(assignmentUser.ID), reqID)
	return nil
}

func (s *AssignmentUserServiceImpl) DeleteAssignmentUser(ctx context.Context, assignmentUserID int) errors.Response {
	reqID, _ := assignmentContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>assignmentUser: delete assignment user started for assignment user id "+strconv.Itoa(assignmentUserID), reqID)

	errResp := s.assignmentUserRepo.DeleteAssignmentUser(ctx, assignmentUserID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return errResp
	}

	log.Info("core>service>assignmentUser: delete assignment user completed for assignment user id "+strconv.Itoa(assignmentUserID), reqID)
	return nil
}

func (s *AssignmentUserServiceImpl) GetAssignmentUserbyID(ctx context.Context, assignmentUserID int) (entity.AssignmentUser, errors.Response) {
	reqID, _ := assignmentContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>assignmentUser: get assignment user started for assignment user id "+strconv.Itoa(assignmentUserID), reqID)

	assignmentUser, errResp := s.assignmentUserRepo.GetAssignmentUserbyID(ctx, assignmentUserID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return entity.AssignmentUser{}, errResp
	}

	log.Info("core>service>assignmentUser: assignment user fetched successfully for assignment user id "+strconv.Itoa(assignmentUserID), reqID)
	return assignmentUser, nil
}

func (s *AssignmentUserServiceImpl) ListAssignmentUser(ctx context.Context) (int, []entity.AssignmentUser, errors.Response) {
	reqID, _ := assignmentContext.GetRequestIDFromContext(ctx)
	log.Info("core>service>assignmentUser: assignment user list started", reqID)

	totalRecords, assignmentUsersEntity, errResp := s.assignmentUserRepo.ListAssignmentUser(ctx)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		return 0, nil, errResp
	}

	log.Info("core>service>assignmentUser: assignment user list completed", reqID)
	return totalRecords, assignmentUsersEntity, nil
}
