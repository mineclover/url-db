package services

import (
	"context"
	"database/sql"
	"log"
	"time"

	"url-db/internal/models"
)

type AttributeRepository interface {
	Create(ctx context.Context, attribute *models.Attribute) error
	GetByID(ctx context.Context, id int) (*models.Attribute, error)
	GetByDomainAndName(ctx context.Context, domainID int, name string) (*models.Attribute, error)
	ListByDomain(ctx context.Context, domainID int) ([]*models.Attribute, error)
	Update(ctx context.Context, attribute *models.Attribute) error
	Delete(ctx context.Context, id int) error
	ExistsByDomainAndName(ctx context.Context, domainID int, name string) (bool, error)
}

type attributeService struct {
	attributeRepo AttributeRepository
	domainRepo    DomainRepository
	logger        *log.Logger
}

func NewAttributeService(attributeRepo AttributeRepository, domainRepo DomainRepository, logger *log.Logger) AttributeService {
	return &attributeService{
		attributeRepo: attributeRepo,
		domainRepo:    domainRepo,
		logger:        logger,
	}
}

func (s *attributeService) CreateAttribute(ctx context.Context, domainID int, req *models.CreateAttributeRequest) (*models.Attribute, error) {
	if err := validatePositiveInteger(domainID, "domainID"); err != nil {
		return nil, err
	}

	req.Name = normalizeString(req.Name)
	req.Type = models.AttributeType(normalizeString(string(req.Type)))
	req.Description = normalizeString(req.Description)

	if err := validateAttributeName(req.Name); err != nil {
		return nil, err
	}

	if err := validateAttributeType(string(req.Type)); err != nil {
		return nil, err
	}

	if err := validateDescription(req.Description); err != nil {
		return nil, err
	}

	domain, err := s.domainRepo.GetByID(ctx, domainID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewDomainNotFoundError(domainID)
		}
		s.logger.Printf("Failed to get domain: %v", err)
		return nil, err
	}

	exists, err := s.attributeRepo.ExistsByDomainAndName(ctx, domainID, req.Name)
	if err != nil {
		s.logger.Printf("Failed to check attribute existence: %v", err)
		return nil, err
	}
	if exists {
		return nil, NewAttributeAlreadyExistsError(domainID, req.Name)
	}

	attribute := &models.Attribute{
		DomainID:    domainID,
		Name:        req.Name,
		Type:        req.Type,
		Description: req.Description,
		CreatedAt:   time.Now(),
	}

	if err := s.attributeRepo.Create(ctx, attribute); err != nil {
		s.logger.Printf("Failed to create attribute: %v", err)
		return nil, err
	}

	s.logger.Printf("Created attribute: %s in domain %s (ID: %d)", attribute.Name, domain.Name, attribute.ID)
	return attribute, nil
}

func (s *attributeService) GetAttribute(ctx context.Context, id int) (*models.Attribute, error) {
	if err := validatePositiveInteger(id, "id"); err != nil {
		return nil, err
	}

	attribute, err := s.attributeRepo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewAttributeNotFoundError(id)
		}
		s.logger.Printf("Failed to get attribute: %v", err)
		return nil, err
	}

	return attribute, nil
}

func (s *attributeService) ListAttributesByDomain(ctx context.Context, domainID int) (*models.AttributeListResponse, error) {
	if err := validatePositiveInteger(domainID, "domainID"); err != nil {
		return nil, err
	}

	domain, err := s.domainRepo.GetByID(ctx, domainID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewDomainNotFoundError(domainID)
		}
		s.logger.Printf("Failed to get domain: %v", err)
		return nil, err
	}

	attributes, err := s.attributeRepo.ListByDomain(ctx, domainID)
	if err != nil {
		s.logger.Printf("Failed to list attributes: %v", err)
		return nil, err
	}

	attributeList := make([]models.Attribute, len(attributes))
	for i, attribute := range attributes {
		attributeList[i] = *attribute
	}

	return &models.AttributeListResponse{
		Attributes: attributeList,
		Domain:     domain,
	}, nil
}

func (s *attributeService) UpdateAttribute(ctx context.Context, id int, req *models.UpdateAttributeRequest) (*models.Attribute, error) {
	if err := validatePositiveInteger(id, "id"); err != nil {
		return nil, err
	}

	req.Description = normalizeString(req.Description)

	if err := validateDescription(req.Description); err != nil {
		return nil, err
	}

	attribute, err := s.attributeRepo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewAttributeNotFoundError(id)
		}
		s.logger.Printf("Failed to get attribute: %v", err)
		return nil, err
	}

	attribute.Description = req.Description

	if err := s.attributeRepo.Update(ctx, attribute); err != nil {
		s.logger.Printf("Failed to update attribute: %v", err)
		return nil, err
	}

	s.logger.Printf("Updated attribute: %s (ID: %d)", attribute.Name, attribute.ID)
	return attribute, nil
}

func (s *attributeService) DeleteAttribute(ctx context.Context, id int) error {
	if err := validatePositiveInteger(id, "id"); err != nil {
		return err
	}

	err := s.attributeRepo.Delete(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return NewAttributeNotFoundError(id)
		}
		s.logger.Printf("Failed to delete attribute: %v", err)
		return err
	}

	s.logger.Printf("Deleted attribute with ID: %d", id)
	return nil
}

func (s *attributeService) ValidateAttributeValue(ctx context.Context, attributeID int, value string) error {
	if err := validatePositiveInteger(attributeID, "attributeID"); err != nil {
		return err
	}

	value = normalizeString(value)

	attribute, err := s.attributeRepo.GetByID(ctx, attributeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return NewAttributeNotFoundError(attributeID)
		}
		s.logger.Printf("Failed to get attribute: %v", err)
		return err
	}

	if err := validateAttributeValue(string(attribute.Type), value); err != nil {
		return NewAttributeValueInvalidError(attributeID, value, err.Error())
	}

	return nil
}
