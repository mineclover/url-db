package services

import "fmt"

type ServiceError struct {
	Code    string
	Message string
	Details interface{}
}

func (e *ServiceError) Error() string {
	return e.Message
}

func NewValidationError(field, message string) *ServiceError {
	return &ServiceError{
		Code:    "VALIDATION_ERROR",
		Message: fmt.Sprintf("Validation failed for field '%s': %s", field, message),
		Details: map[string]string{"field": field, "message": message},
	}
}

func NewDomainAlreadyExistsError(name string) *ServiceError {
	return &ServiceError{
		Code:    "DOMAIN_ALREADY_EXISTS",
		Message: fmt.Sprintf("Domain '%s' already exists", name),
		Details: map[string]string{"domain": name},
	}
}

func NewNodeAlreadyExistsError(url string) *ServiceError {
	return &ServiceError{
		Code:    "NODE_ALREADY_EXISTS",
		Message: fmt.Sprintf("Node with URL '%s' already exists", url),
		Details: map[string]string{"url": url},
	}
}

func NewAttributeAlreadyExistsError(domainID int, name string) *ServiceError {
	return &ServiceError{
		Code:    "ATTRIBUTE_ALREADY_EXISTS",
		Message: fmt.Sprintf("Attribute '%s' already exists in domain %d", name, domainID),
		Details: map[string]interface{}{"domain_id": domainID, "name": name},
	}
}

func NewInvalidCompositeKeyError(key, reason string) *ServiceError {
	return &ServiceError{
		Code:    "INVALID_COMPOSITE_KEY",
		Message: fmt.Sprintf("Invalid composite key '%s': %s", key, reason),
		Details: map[string]string{"key": key, "reason": reason},
	}
}

func NewDomainNotFoundError(id int) *ServiceError {
	return &ServiceError{
		Code:    "DOMAIN_NOT_FOUND",
		Message: fmt.Sprintf("Domain with id %d not found", id),
		Details: map[string]int{"id": id},
	}
}

func NewNodeNotFoundError(id int) *ServiceError {
	return &ServiceError{
		Code:    "NODE_NOT_FOUND",
		Message: fmt.Sprintf("Node with id %d not found", id),
		Details: map[string]int{"id": id},
	}
}

func NewAttributeNotFoundError(id int) *ServiceError {
	return &ServiceError{
		Code:    "ATTRIBUTE_NOT_FOUND",
		Message: fmt.Sprintf("Attribute with id %d not found", id),
		Details: map[string]int{"id": id},
	}
}

func NewNodeAttributeNotFoundError(id int) *ServiceError {
	return &ServiceError{
		Code:    "NODE_ATTRIBUTE_NOT_FOUND",
		Message: fmt.Sprintf("Node attribute with id %d not found", id),
		Details: map[string]int{"id": id},
	}
}

func NewAttributeValueInvalidError(attributeID int, value, reason string) *ServiceError {
	return &ServiceError{
		Code:    "ATTRIBUTE_VALUE_INVALID",
		Message: fmt.Sprintf("Invalid value '%s' for attribute %d: %s", value, attributeID, reason),
		Details: map[string]interface{}{"attribute_id": attributeID, "value": value, "reason": reason},
	}
}

func NewBusinessLogicError(message string) *ServiceError {
	return &ServiceError{
		Code:    "BUSINESS_LOGIC_ERROR",
		Message: message,
		Details: nil,
	}
}
