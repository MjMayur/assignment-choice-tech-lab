package apimodel

type Response struct {
	StatusCode int         `json:"code"`
	Status     string      `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

type InvalidValidationError struct {
	Field string `json:"Field"`
	Msg   string `json:"Msg"`
}

type ValidationErrorResponse struct {
	StatusCode int                      `json:"code"`
	Status     string                   `json:"status"`
	Message    []InvalidValidationError `json:"message"`
	Data       interface{}              `json:"data,omitempty"`
}
