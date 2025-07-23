package services

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"url-db/internal/models"
)

type compositeKeyService struct {
	toolName string
}

func NewCompositeKeyService(toolName string) CompositeKeyService {
	return &compositeKeyService{
		toolName: toolName,
	}
}

func (s *compositeKeyService) Create(domainName string, id int) string {
	normalizedDomain := s.normalizeName(domainName)
	return fmt.Sprintf("%s:%s:%d", s.toolName, normalizedDomain, id)
}

func (s *compositeKeyService) Parse(compositeKey string) (*models.CompositeKey, error) {
	if compositeKey == "" {
		return nil, NewInvalidCompositeKeyError(compositeKey, "composite key is empty")
	}

	parts := strings.Split(compositeKey, ":")
	if len(parts) != 3 {
		return nil, NewInvalidCompositeKeyError(compositeKey, "invalid format, expected 'tool:domain:id'")
	}

	toolName := parts[0]
	domainName := parts[1]
	idStr := parts[2]

	if toolName == "" {
		return nil, NewInvalidCompositeKeyError(compositeKey, "tool name is empty")
	}

	if toolName != s.toolName {
		return nil, NewInvalidCompositeKeyError(compositeKey, fmt.Sprintf("invalid tool name, expected '%s'", s.toolName))
	}

	if domainName == "" {
		return nil, NewInvalidCompositeKeyError(compositeKey, "domain name is empty")
	}

	if idStr == "" {
		return nil, NewInvalidCompositeKeyError(compositeKey, "ID is empty")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, NewInvalidCompositeKeyError(compositeKey, "ID must be a valid integer")
	}

	if id <= 0 {
		return nil, NewInvalidCompositeKeyError(compositeKey, "ID must be positive")
	}

	return &models.CompositeKey{
		ToolName:   toolName,
		DomainName: domainName,
		ID:         id,
	}, nil
}

func (s *compositeKeyService) Validate(compositeKey string) error {
	_, err := s.Parse(compositeKey)
	return err
}

func (s *compositeKeyService) GetToolName() string {
	return s.toolName
}

func (s *compositeKeyService) normalizeName(name string) string {
	normalized := strings.ToLower(name)

	reg := regexp.MustCompile(`[^a-z0-9\-_]`)
	normalized = reg.ReplaceAllString(normalized, "-")

	reg = regexp.MustCompile(`-+`)
	normalized = reg.ReplaceAllString(normalized, "-")

	normalized = strings.Trim(normalized, "-")

	if normalized == "" {
		normalized = "default"
	}

	return normalized
}

func (s *compositeKeyService) DenormalizeName(normalizedName string) string {
	return strings.ReplaceAll(normalizedName, "-", " ")
}

func (s *compositeKeyService) ExtractComponents(compositeKey string) (toolName, domainName string, id int, err error) {
	ck, err := s.Parse(compositeKey)
	if err != nil {
		return "", "", 0, err
	}

	return ck.ToolName, ck.DomainName, ck.ID, nil
}

func (s *compositeKeyService) IsValid(compositeKey string) bool {
	return s.Validate(compositeKey) == nil
}

func (s *compositeKeyService) GenerateFromNode(node *models.Node, domainName string) string {
	return s.Create(domainName, node.ID)
}

func (s *compositeKeyService) GenerateFromDomain(domain *models.Domain) string {
	return s.Create(domain.Name, domain.ID)
}
