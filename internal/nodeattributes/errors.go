package nodeattributes

import "errors"

var (
	ErrNodeAttributeNotFound       = errors.New("node attribute not found")
	ErrNodeNotFound                = errors.New("node not found")
	ErrAttributeNotFound           = errors.New("attribute not found")
	ErrNodeAttributeValueInvalid   = errors.New("node attribute value invalid")
	ErrNodeAttributeOrderInvalid   = errors.New("node attribute order invalid")
	ErrNodeAttributeDomainMismatch = errors.New("node and attribute domain mismatch")
	ErrNodeAttributeExists         = errors.New("node attribute already exists")
	ErrInvalidAttributeType        = errors.New("invalid attribute type")
	ErrOrderIndexRequired          = errors.New("order index required for ordered_tag type")
	ErrOrderIndexNotAllowed        = errors.New("order index not allowed for this attribute type")
	ErrInvalidOrderIndex           = errors.New("invalid order index")
	ErrDuplicateOrderIndex         = errors.New("duplicate order index")
)
