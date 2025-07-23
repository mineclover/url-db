package attribute

import "errors"

var (
	ErrDomainNotFound    = errors.New("domain not found")
	ErrAttributeNotFound = errors.New("attribute not found")
	ErrInvalidRequest    = errors.New("invalid request")
)