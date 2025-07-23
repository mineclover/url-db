package usecases

import "url-db/internal/models"

// AttributeType aliases for internal use
type AttributeType = models.AttributeType

const (
	AttributeTypeTag        = models.AttributeTypeTag
	AttributeTypeOrderedTag = models.AttributeTypeOrderedTag
	AttributeTypeNumber     = models.AttributeTypeNumber
	AttributeTypeString     = models.AttributeTypeString
	AttributeTypeMarkdown   = models.AttributeTypeMarkdown
	AttributeTypeImage      = models.AttributeTypeImage
)

// IsValidAttributeType checks if the given attribute type is valid
func IsValidAttributeType(t AttributeType) bool {
	switch t {
	case AttributeTypeTag, AttributeTypeOrderedTag, AttributeTypeNumber,
		AttributeTypeString, AttributeTypeMarkdown, AttributeTypeImage:
		return true
	default:
		return false
	}
}

// GetSupportedAttributeTypes returns a list of all supported attribute types
func GetSupportedAttributeTypes() []AttributeType {
	return []AttributeType{
		AttributeTypeTag,
		AttributeTypeOrderedTag,
		AttributeTypeNumber,
		AttributeTypeString,
		AttributeTypeMarkdown,
		AttributeTypeImage,
	}
}
