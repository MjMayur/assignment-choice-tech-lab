package validation

import (
	"fmt"
	"io"
	"strconv"

	apiModel "project/apimodel"
	"project/pkg/errors"

	"github.com/gin-gonic/gin"
	strcase "github.com/iancoleman/strcase"
	"github.com/rs/zerolog/log"
)

func DecodeAndValidate(r io.Reader, requestInstance interface{}, c *gin.Context) ([]apiModel.InvalidValidationError, errors.Response) {
	// decode the request
	InvalidValidationErrors := []apiModel.InvalidValidationError{}
	err := c.ShouldBindJSON(requestInstance)
	if err != nil {
		fields, customValidationErrs := ValidationError(err)
		for i := 0; i < len(fields); i++ {
			log.Error().Str("validation  ", "Error").Msg(customValidationErrs[i])
			InvalidValidationErrors = append(InvalidValidationErrors, apiModel.InvalidValidationError{
				Field: strcase.ToLowerCamel(fields[i]),
				Msg:   customValidationErrs[i],
			})
		}

		return InvalidValidationErrors, errors.ResponseInternalServerError(err.Error())
	}

	return nil, nil
}

func DecodeAndValidateForQueryParams(c *gin.Context, requestInstance interface{}) errors.Response {
	// decode the request
	err := c.ShouldBindQuery(requestInstance)
	if err != nil {
		if numErr, ok := err.(*strconv.NumError); ok {
			msg1 := fmt.Sprintf("Invalid numeric value is " + numErr.Num)
			log.Error().Str("validation  ", "Error").Msg(err.Error())
			return errors.ResponseBadRequestError(msg1)
		}

		customValidationErrs := ValidationQueryParamsError(err)
		log.Error().Str("validation  ", "Error").Msg(err.Error())
		return errors.ResponseBadRequestError(customValidationErrs)
	}

	return nil
}

func DecodeAndValidateForm(c *gin.Context, requestInstance interface{}) ([]apiModel.InvalidValidationError, errors.Response) {

	InvalidValidationErrors := []apiModel.InvalidValidationError{}
	err := c.ShouldBind(requestInstance)
	if err != nil {
		fields, customValidationErrs := ValidationError(err)
		for i := 0; i < len(fields); i++ {
			log.Error().Str("validation  ", "Error").Msg(customValidationErrs[i])
			InvalidValidationErrors = append(InvalidValidationErrors, apiModel.InvalidValidationError{
				Field: fields[i],
				Msg:   customValidationErrs[i],
			})
		}

		return InvalidValidationErrors, errors.ResponseInternalServerError(err.Error())
	}

	return nil, nil
}
