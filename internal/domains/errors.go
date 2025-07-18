package domains

import "errors"

var (
	ErrDomainNotFound        = errors.New("domain not found")
	ErrDomainAlreadyExists   = errors.New("domain already exists")
	ErrDomainNameInvalid     = errors.New("domain name is invalid")
	ErrDomainHasDependencies = errors.New("domain has dependencies and cannot be deleted")
	ErrValidation            = errors.New("validation error")
)

const (
	ErrorCodeDomainNotFound        = "DOMAIN_NOT_FOUND"
	ErrorCodeDomainAlreadyExists   = "DOMAIN_ALREADY_EXISTS"
	ErrorCodeDomainNameInvalid     = "DOMAIN_NAME_INVALID"
	ErrorCodeDomainHasDependencies = "DOMAIN_HAS_DEPENDENCIES"
	ErrorCodeValidation            = "VALIDATION_ERROR"
)

type DomainError struct {
	Code    string
	Message string
	Err     error
}

func (e *DomainError) Error() string {
	return e.Message
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

func NewDomainError(code, message string, err error) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func NewValidationError(message string) *DomainError {
	return &DomainError{
		Code:    ErrorCodeValidation,
		Message: message,
		Err:     ErrValidation,
	}
}

func NewDomainNotFoundError(id int) *DomainError {
	return &DomainError{
		Code:    ErrorCodeDomainNotFound,
		Message: "Domain not found",
		Err:     ErrDomainNotFound,
	}
}

func NewDomainAlreadyExistsError(name string) *DomainError {
	return &DomainError{
		Code:    ErrorCodeDomainAlreadyExists,
		Message: "Domain already exists",
		Err:     ErrDomainAlreadyExists,
	}
}
