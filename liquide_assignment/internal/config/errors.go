package config

import "net/http"

type Errors struct {
	ActualError string    `json:"-"`
	Message     string    `json:"message,omitempty"`
	Status      Status    `json:"status,omitempty"`
	Code        ErrorCode `json:"code,omitempty"`
	StatusCode  int       `json:"-"`
}

func (e Errors) ErrorMessage() string {
	return e.Message
}

func (e Errors) Error() string {
	return e.ActualError
}

func (e Errors) ErrorCode() ErrorCode {
	return e.Code
}

func WrapWithStatus(err error, ee Errors, status int) Errors {
	ee.ActualError = err.Error()
	ee.StatusCode = status
	return ee
}

func Wrap(err error, ee Errors) Errors {
	ee.ActualError = err.Error()
	return ee
}

type ErrorCode string

/*
Error codes
1xx -> System Error
2xx ->
3xx -> Unauthorized/Invalid Request
4xx -> Redis Error
5xx -> Mongo Error
*/

const (
	// System Error Codes
	PasswordHashError  ErrorCode = "LIQ-100"
	TooManyRequests    ErrorCode = "LIQ-101"
	SessionIssueError  ErrorCode = "LIQ-102"
	SessionDecodeError ErrorCode = "LIQ-103"

	// 2xx Error Codes

	// Unauthorized/Invalid Request Error Codes
	InvalidRequest    ErrorCode = "LIQ-300"
	InvalidUser       ErrorCode = "LIQ-301"
	InvalidPassword   ErrorCode = "LIQ-302"
	UserAlreadyExists ErrorCode = "LIQ-303"
	ValidationError   ErrorCode = "LIQ-304"
	Unauthorized      ErrorCode = "LIQ-305"

	// Redis Error Codes
	DatabaseUpdateErrorRedis ErrorCode = "LIQ-400"
	InternalServerErrorRedis ErrorCode = "LIQ-401"

	// Mongo Error Codes
	DatabaseUpdateErrorMongo ErrorCode = "LIQ-500"
	InternalServerErrorMongo ErrorCode = "LIQ-501"
)

var (
	ErrTooManyRequests    = Errors{Message: "Too many requests", Code: TooManyRequests, Status: Failure, StatusCode: http.StatusTooManyRequests}
	ErrPasswordHashError  = Errors{Message: "Failed to hash password", Code: PasswordHashError, Status: Failure, StatusCode: http.StatusInternalServerError}
	ErrSessionIssueError  = Errors{Message: "Failed to issue session", Code: SessionIssueError, Status: Failure, StatusCode: http.StatusInternalServerError}
	ErrSomethingWentWrong = Errors{Message: "Something went wrong", Code: SessionDecodeError, Status: Failure, StatusCode: http.StatusInternalServerError}

	ErrInvalidRequest    = Errors{Message: "Invalid request", Code: InvalidRequest, Status: Failure, StatusCode: http.StatusBadRequest}
	ErrInvalidUser       = Errors{Message: "Invalid user", Code: InvalidUser, Status: Failure, StatusCode: http.StatusBadRequest}
	ErrInvalidPassword   = Errors{Message: "Invalid password", Code: InvalidPassword, Status: Failure, StatusCode: http.StatusBadRequest}
	ErrUserAlreadyExists = Errors{Message: "User already exists", Code: UserAlreadyExists, Status: Failure, StatusCode: http.StatusBadRequest}
	ErrValidationError   = Errors{Message: "Validation error", Code: ValidationError, Status: Failure, StatusCode: http.StatusBadRequest}
	ErrUnauthorized      = Errors{Message: "Unauthorized", Code: Unauthorized, Status: Failure, StatusCode: http.StatusUnauthorized}

	ErrDatabaseUpdateErrorRedis = Errors{Message: "Failed to update database", Code: DatabaseUpdateErrorRedis, Status: Failure, StatusCode: http.StatusInternalServerError}
	ErrInternalServerErrorRedis = Errors{Message: "Something went wrong", Code: InternalServerErrorRedis, Status: Failure, StatusCode: http.StatusInternalServerError}

	ErrDatabaseUpdateErrorMongo = Errors{Message: "Failed to update database", Code: DatabaseUpdateErrorMongo, Status: Failure, StatusCode: http.StatusInternalServerError}
	ErrInternalServerErrorMongo = Errors{Message: "Something went wrong", Code: InternalServerErrorMongo, Status: Failure, StatusCode: http.StatusInternalServerError}
)
