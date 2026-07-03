package utils

import (
	"encoding/json"

	apiModel "project/apimodel"
	"project/pkg/errors"

	"github.com/gin-gonic/gin"
)

const (
	STATUS_SUCCESS               = "SUCCESS"
	STATUS_FAILED                = "FAILED"
	STATUS_INTERNAL_SERVER_ERROR = "INTENAL SERVER ERROR"
)

func DataResponse(c *gin.Context, statusCode int, msg string, payload interface{}) {
	res := apiModel.Response{
		StatusCode: statusCode,
		Status:     STATUS_SUCCESS,
		Message:    msg,
		Data:       payload,
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(res)
	c.Writer.WriteHeader(statusCode)
	c.Writer.Write(response)
}

func ErrorResponse(c *gin.Context, errResp errors.Response, payload interface{}) {

	if errResp == nil {
		ErrorResponse(c, errors.ResponseInternalServerError(STATUS_INTERNAL_SERVER_ERROR), nil)
		return
	}

	errMsg := errResp.Error()

	res := apiModel.Response{
		StatusCode: errResp.StatusCode(),
		Status:     STATUS_FAILED,
		Message:    errMsg,
		Data:       payload,
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(res)
	c.Writer.WriteHeader(errResp.StatusCode())
	c.Writer.Write(response)
}

func ValidationErrorResponse(c *gin.Context, fieldErrors []apiModel.InvalidValidationError, payload interface{}) {

	if fieldErrors == nil {
		ErrorResponse(c, errors.ResponseInternalServerError(STATUS_INTERNAL_SERVER_ERROR), nil)
		return
	}

	resp := apiModel.ValidationErrorResponse{
		StatusCode: 400,
		Status:     STATUS_FAILED,
		Message:    fieldErrors,
		Data:       payload,
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(resp)
	c.Writer.WriteHeader(400)
	c.Writer.Write(response)
}
