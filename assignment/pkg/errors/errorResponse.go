package errors

import "net/http"

type Response interface {
	Error() string
	StatusCode() int
}

type ResponseImpl struct {
	statusCode int
	msg        string
}

func (h ResponseImpl) Error() string {
	return h.msg
}

func (h ResponseImpl) StatusCode() int {
	return h.statusCode
}

// generator function for the general http error
func NewResponseError(statusCode int, msg string) Response {
	return ResponseImpl{
		statusCode: statusCode,
		msg:        msg,
	}
}

// The “HTTP 500 Internal Server Error“ server error response code indicates that the server encountered an unexpected condition.
//
// - that prevented it from fulfilling the request.
func ResponseInternalServerError(msg string) Response {
	return NewResponseError(http.StatusInternalServerError, msg)
}

// The “HTTP 400 Bad Request“ response status code indicates that the server cannot or will not process the request.
//
// - due to something that is perceived to be a client error.
func ResponseBadRequestError(msg string) Response {
	return NewResponseError(http.StatusBadRequest, msg)
}

// The “HTTP 403 Forbidden“ response status code indicates that the server understands the request but refuses to authorize it.
//
// - Access tied to application logic, such as insufficient rights to a resource.
func ResponseForbiddenRequestError(msg string) Response {
	return NewResponseError(http.StatusForbidden, msg)
}

// The “HTTP 404 Not Found“ response status code indicates that the server cannot find the requested resource.
//
// - often occurs due to broken or dead links
func ResponseNotFoundError(msg string) Response {
	return NewResponseError(http.StatusNotFound, msg)
}

// The “401 Unauthorized“ response status code indicates that the client request has not been completed.
//
// - because it lacks valid authentication credentials for the requested resource.
func ResponseUnauthorizedError(msg string) Response {
	return NewResponseError(http.StatusUnauthorized, msg)
}

// The “HTTP 417 Expectation Failed“ client error response code indicates that the expectation given in the request's Expect header could not be met.
func ResponseExpectionFailedError(msg string) Response {
	return NewResponseError(http.StatusExpectationFailed, msg)
}

// The “HTTP 409 Conflict“ response status code indicate a request conflict with the current state of the target resource.
//
// - Conflicts are most likely to occur in response to a `PUT` request. So 409 Conflict
func ResponseConflictError(msg string) Response {
	return NewResponseError(http.StatusConflict, msg)
}

// The HTTP 204 No Content success status response code indicates that a request has succeeded,
//
// but that the client doesn't need to navigate away from its current page
func ResponseNoContentError(msg string) Response {
	return NewResponseError(http.StatusNoContent, msg)
}

func ResponseLoginTimeoutError(msg string) Response {
	return NewResponseError(440, msg)
}
