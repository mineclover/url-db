package nodes

import "errors"

var (
	ErrNodeNotFound       = errors.New("node not found")
	ErrNodeAlreadyExists  = errors.New("node already exists")
	ErrNodeURLInvalid     = errors.New("node url invalid")
	ErrNodeDomainNotFound = errors.New("node domain not found")
	ErrNodeHasAttributes  = errors.New("node has attributes")
	ErrValidationError    = errors.New("validation error")
	ErrConflict           = errors.New("conflict")
)
