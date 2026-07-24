package errors

type AppError struct {
	Code       string
	Message    string
	StatusCode int
	Details    any
}

func (e *AppError) Error() string {
	return e.Message
}
