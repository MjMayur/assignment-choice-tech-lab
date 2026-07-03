package validation

import (
	"regexp"

	apiModel "project/apimodel"
)

func ValidatePhone(phoneNumber string) *apiModel.InvalidValidationError {
	regularExpression := regexp.MustCompile(`^[7896]\d{9}$`)
	ok := regularExpression.MatchString(phoneNumber)

	if !ok {
		validationErrors := apiModel.InvalidValidationError{
			Field: "phone",
			Msg:   "invalid phone number",
		}

		return &validationErrors
	}

	return nil
}
