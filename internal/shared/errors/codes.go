package errors

const (
	// Validation
	CodeValidation = "VALIDATION_ERROR"

	// Authentication
	CodeUnauthorized = "UNAUTHORIZED"
	CodeForbidden    = "FORBIDDEN"

	// Business
	CodeNotFound = "NOT_FOUND"
	CodeConflict = "CONFLICT"

	// Internal
	CodeInternal = "INTERNAL_ERROR"
)
