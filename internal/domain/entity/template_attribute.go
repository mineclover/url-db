package entity

import (
	"fmt"
	"time"
)

// TemplateAttribute represents a template attribute association
type TemplateAttribute struct {
	id          int
	templateID  int
	attributeID int
	value       string
	orderIndex  *int
	createdAt   time.Time
}

// NewTemplateAttribute creates a new template attribute
func NewTemplateAttribute(templateID, attributeID int, value string, orderIndex *int) (*TemplateAttribute, error) {
	if templateID <= 0 {
		return nil, fmt.Errorf("template ID must be positive, got %d", templateID)
	}
	if attributeID <= 0 {
		return nil, fmt.Errorf("attribute ID must be positive, got %d", attributeID)
	}
	if value == "" {
		return nil, fmt.Errorf("attribute value cannot be empty")
	}

	return &TemplateAttribute{
		templateID:  templateID,
		attributeID: attributeID,
		value:       value,
		orderIndex:  orderIndex,
		createdAt:   time.Now(),
	}, nil
}

// ID returns the template attribute ID
func (ta *TemplateAttribute) ID() int {
	return ta.id
}

// TemplateID returns the template ID
func (ta *TemplateAttribute) TemplateID() int {
	return ta.templateID
}

// AttributeID returns the attribute ID
func (ta *TemplateAttribute) AttributeID() int {
	return ta.attributeID
}

// Value returns the attribute value
func (ta *TemplateAttribute) Value() string {
	return ta.value
}

// OrderIndex returns the order index for ordered attributes
func (ta *TemplateAttribute) OrderIndex() *int {
	return ta.orderIndex
}

// CreatedAt returns the creation timestamp
func (ta *TemplateAttribute) CreatedAt() time.Time {
	return ta.createdAt
}

// SetID sets the template attribute ID (used by repository)
func (ta *TemplateAttribute) SetID(id int) {
	ta.id = id
}

// UpdateValue updates the attribute value
func (ta *TemplateAttribute) UpdateValue(value string) error {
	if value == "" {
		return fmt.Errorf("attribute value cannot be empty")
	}
	ta.value = value
	return nil
}

// UpdateOrderIndex updates the order index
func (ta *TemplateAttribute) UpdateOrderIndex(orderIndex *int) {
	ta.orderIndex = orderIndex
}

// SetCreatedAt sets the creation timestamp (used by repository)
func (ta *TemplateAttribute) SetCreatedAt(createdAt time.Time) {
	ta.createdAt = createdAt
}

// String returns a string representation of the template attribute
func (ta *TemplateAttribute) String() string {
	orderStr := "nil"
	if ta.orderIndex != nil {
		orderStr = fmt.Sprintf("%d", *ta.orderIndex)
	}
	return fmt.Sprintf("TemplateAttribute{ID: %d, TemplateID: %d, AttributeID: %d, Value: %s, OrderIndex: %s}",
		ta.id, ta.templateID, ta.attributeID, ta.value, orderStr)
}

// TemplateAttributeWithDetails combines template attribute with attribute details
type TemplateAttributeWithDetails struct {
	TemplateAttribute *TemplateAttribute
	AttributeName     string
	AttributeType     string
	AttributeDesc     string
}

// NewTemplateAttributeWithDetails creates a new template attribute with details
func NewTemplateAttributeWithDetails(ta *TemplateAttribute, name, attrType, desc string) *TemplateAttributeWithDetails {
	return &TemplateAttributeWithDetails{
		TemplateAttribute: ta,
		AttributeName:     name,
		AttributeType:     attrType,
		AttributeDesc:     desc,
	}
}
