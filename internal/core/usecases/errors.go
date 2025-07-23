package attributes

import "errors"

var (
	// Attribute errors
	ErrAttributeAlreadyExists = errors.New("attribute already exists")
	ErrAttributeNotFound      = errors.New("attribute not found")
	ErrAttributeTypeInvalid   = errors.New("attribute type invalid")
	ErrAttributeHasValues     = errors.New("attribute has values")
	ErrDomainNotFound         = errors.New("domain not found")

	// Validation errors
	ErrAttributeNameRequired = errors.New("attribute name is required")
	ErrAttributeNameTooLong  = errors.New("attribute name too long")
	ErrAttributeTypeRequired = errors.New("attribute type is required")
	ErrDescriptionTooLong    = errors.New("description too long")

	// Value validation errors
	ErrValueRequired      = errors.New("value is required")
	ErrValueTooLong       = errors.New("value too long")
	ErrValueInvalid       = errors.New("value invalid")
	ErrInvalidURL         = errors.New("invalid URL")
	ErrInvalidNumber      = errors.New("invalid number")
	ErrInvalidMarkdown    = errors.New("invalid markdown")
	ErrOrderIndexRequired = errors.New("order index required for ordered tag")
)
