package apperrors

type APIError struct {
	Status  int
	Code    string
	Message string
	Details any
}

// Membuat custom API error
func New(status int, code string, message string, details any) *APIError {
	return &APIError{
		Status:  status,
		Code:    code,
		Message: message,
		Details: details,
	}
}

// Implementasi interface error milik Go
func (e *APIError) Error() string {
	return e.Message
}
