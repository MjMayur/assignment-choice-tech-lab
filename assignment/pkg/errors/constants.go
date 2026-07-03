package errors

var DYNAMIC_ERROR string

const (
	BAD_REQUEST_ERROR         = "Bad Request"           // HTTP code 400
	UNAUTHORIZED_ERROR        = "Unauthorized"          // HTTP code 401
	FORBIDDEN_ERROR           = "Forbidden"             // HTTP code 403
	NOT_FOUND_ERROR           = "Not Found"             // HTTP code 404
	METHOD_NOT_ALLOWED_ERROR  = "Method Not Allowed"    // HTTP code 405
	REQUEST_TIMEOUT_ERROR     = "Request Timeout"       // HTTP code 408
	CONFLICT_ERROR            = "Conflict"              // HTTP code 409
	EXPECTATION_FAILED_ERROR  = "Expectation Failed"    // HTTP code 417
	INTERNAL_SERVER_ERROR     = "Internal server error" // HTTP code 500
	BAD_GATEWAY_ERROR         = "Bad Gateway"           // HTTP code 502
	SERVICE_UNAVAILABLE_ERROR = "Service Unavailable"   // HTTP code 503
	GATEWAY_TIMEOUT_ERROR     = "Gateway Timeout"       // HTTP code 504
	NO_CONTENT_ERROR          = "No Content"            // HTTP code 204
)

const (
	COULD_NOT_CREATE_SESSION = "Could not create the session."
	COULD_NOT_GET_SESSION    = "Could not get the session."
	DELETE_SESSION_ERR       = "Cannot delete the session."
	GET_SESSION_FAILED       = "Failed to get the session."
	UNAUTHORIZED_USER        = "Unauthorized user."
	MSG_ACCESS_DENIED        = "Access denied."
	SESSION_ERROR            = "Session Timeout"
)
