package sys

// AppError defines a structured error for the frontend
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NewError creates a new AppError
func NewError(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}
