package web

import (
	"net/http"
	"path/filepath"

	coreAPIModel "project/apimodel/core"
	"project/internal/converter"
	"project/pkg/context"
	"project/pkg/errors"
	"project/pkg/log"
	httpUtils "project/pkg/utils"
	"project/pkg/validation"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h CoreHandlerRegistry) UploadCsv(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>user: import daily records from CSV started", reqID)

	// Parse multipart form
	err := c.Request.ParseMultipartForm(10 << 20) // 10MB max
	if err != nil {
		log.Error("Failed to parse multipart form: "+err.Error(), reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("Failed to parse form data"), nil)
		return
	}

	form := c.Request.MultipartForm
	if form == nil {
		log.Error("Multipart form is nil", reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("Invalid form data"), nil)
		return
	}

	files := form.File["csvFile"]

	if len(files) == 0 {
		log.Error("csv file is required", reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("CSV file is required"), nil)
		return
	}

	if len(files) > 1 {
		log.Error("only 1 csv file is allowed", reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("Only 1 CSV file is allowed"), nil)
		return
	}

	csvFile := files[0]

	// Validate file size (max 10MB)
	if csvFile.Size > 10<<20 {
		log.Error("File size exceeds 10MB limit", reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("File size exceeds 10MB limit"), nil)
		return
	}

	filename := filepath.Base(csvFile.Filename)
	ext := strings.ToLower(filepath.Ext(filename))

	if ext != ".csv" {
		log.Error("not a CSV file: "+filename, reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("Please upload a CSV file"), nil)
		return
	}

	// Process the CSV file
	storeMap, ignoreMap, successMsg, errResp := h.Options.AssignmentUserService.StoreCsvRecords(c.Request.Context(), csvFile)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	// Prepare response data
	responseData := map[string]interface{}{
		"successResponse": storeMap,
		"invalidResponse": ignoreMap,
		"message":         successMsg,
	}

	log.Info("core>web>user: import daily records from CSV completed successfully", reqID)
	httpUtils.DataResponse(c, http.StatusOK, "CSV records uploaded successfully", responseData)
}

func (h CoreHandlerRegistry) CreateAssignmentUserHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>assignmentUser: create assignment user started", reqID)

	assignmentUserRequest := &coreAPIModel.AssignmentUser{}
	validationErrors, err := validation.DecodeAndValidate(c.Request.Body, assignmentUserRequest, c)
	if err != nil {
		log.Error(err.Error(), reqID)
		httpUtils.ValidationErrorResponse(c, validationErrors, nil)
		return
	}

	assignmentUserEntity := converter.AssignmentUserAPIToAssignmentUserEntity(assignmentUserRequest)

	assignmentUserID, errResp := h.Options.AssignmentUserService.CreateAssignmentUser(c.Request.Context(), assignmentUserEntity)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	log.Info("core>web>assignmentUser: create assignment user completed & assignment user id is "+strconv.Itoa(assignmentUserID), reqID)
	httpUtils.DataResponse(c, http.StatusOK, "AssignmentUser created successfully", nil)
}

func (h CoreHandlerRegistry) PartialUpdateAssignmentUserHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>assignmentUser: partial update assignment user started", reqID)

	assignmentUserIDStr := c.Param("id")
	assignmentUserIDStr = strings.TrimSpace(assignmentUserIDStr)

	assignmentUserID, err := strconv.Atoi(assignmentUserIDStr)
	if err != nil {
		log.Error("invalid assignmentUser id", reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("invalid assignmentUser id"), nil)
		return
	}

	log.Info("core>web>assignmentUser: partial assignment user update started for assignment user id "+assignmentUserIDStr, reqID)

	assignmentUserRequest := &coreAPIModel.UpdateAssignmentUserRequest{}
	validationErrors, err := validation.DecodeAndValidate(c.Request.Body, assignmentUserRequest, c)
	if err != nil {
		log.Error(err.Error(), reqID)
		httpUtils.ValidationErrorResponse(c, validationErrors, nil)
		return
	}

	assignmentUserEntity := converter.UpdateAssignmentUserAPIRequestToAssignmentUserEntity(assignmentUserRequest)
	assignmentUserEntity.ID = assignmentUserID

	errResp := h.Options.AssignmentUserService.PartialUpdateAssignmentUser(c.Request.Context(), assignmentUserEntity)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	log.Info("core>web>assignmentUser: partial update assignment user completed successfully for assignment user id "+assignmentUserIDStr, reqID)
	httpUtils.DataResponse(c, http.StatusOK, "partial update assignment user completed successfully", nil)
}

func (h CoreHandlerRegistry) DeleteAssignmentUserHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>assignmentUser: delete assignment user started", reqID)

	assignmentUserIDstr := c.Param("id")
	assignmentUserIDstr = strings.TrimSpace(assignmentUserIDstr)

	assignmentUserID, err := strconv.Atoi(assignmentUserIDstr)
	if err != nil {
		log.Error("invalid assignmentUser id", reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("invalid assignmentUser id"), nil)
		return
	}

	log.Info("core>web>assignmentUser: delete assignment user started for assignment user id "+strconv.Itoa(assignmentUserID), reqID)

	errResp := h.Options.AssignmentUserService.DeleteAssignmentUser(c.Request.Context(), assignmentUserID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	log.Info("core>web>assignmentUser: delete assignment user completed for assignment user id "+strconv.Itoa(assignmentUserID), reqID)
	httpUtils.DataResponse(c, http.StatusOK, "assignment user deleted successfully", nil)
}

func (h CoreHandlerRegistry) GetAssignmentUserbyIDHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>assignmentUser: get assignment user started", reqID)

	assignmentUserIDstr := c.Param("id")
	assignmentUserIDstr = strings.TrimSpace(assignmentUserIDstr)

	assignmentUserID, err := strconv.Atoi(assignmentUserIDstr)
	if err != nil {
		log.Error("invalid assignmentUser id", reqID)
		httpUtils.ErrorResponse(c, errors.ResponseBadRequestError("invalid assignmentUser id"), nil)
		return
	}

	log.Info("core>web>assignmentUser:get assignment user started for assignment user id "+assignmentUserIDstr, reqID)

	assignmentUser, errResp := h.Options.AssignmentUserService.GetAssignmentUserbyID(c.Request.Context(), assignmentUserID)
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	assignmentUserResponse := converter.AssignmentUserEntityToAssignmentUserAPIModelResponse(assignmentUser)

	log.Info("core>web>assignmentUser: assignment user fetched successfully for assignment user id "+strconv.Itoa(assignmentUserID), reqID)
	httpUtils.DataResponse(c, http.StatusOK, "assignment user fetched successfully", assignmentUserResponse)
}

func (h CoreHandlerRegistry) ListAssignmentUserHandler(c *gin.Context) {
	reqID, _ := context.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>assignmentUser: assignment user list started", reqID)

	totalRecords, assignmentUsersEntity, errResp := h.Options.AssignmentUserService.ListAssignmentUser(c.Request.Context())
	if errResp != nil {
		log.Error(errResp.Error(), reqID)
		httpUtils.ErrorResponse(c, errResp, nil)
		return
	}

	assignmentUsersListResponse := []coreAPIModel.AssignmentUserResponse{}
	for _, assignmentUserEntity := range assignmentUsersEntity {
		assignmentUserresponse := converter.AssignmentUserEntityToAssignmentUserAPIModelResponse(assignmentUserEntity)
		assignmentUsersListResponse = append(assignmentUsersListResponse, assignmentUserresponse)
	}

	var assignmentUserResponse coreAPIModel.AssignmentUserListResponse
	assignmentUserResponse.TotalRecords = totalRecords
	assignmentUserResponse.NoOfRecordsPerPage = 0
	assignmentUserResponse.AssignmentUser = assignmentUsersListResponse

	log.Info("core>web>assignmentUser: assignment user list completed", reqID)
	httpUtils.DataResponse(c, http.StatusOK, "assignment users fetched successfully", assignmentUserResponse)
}
