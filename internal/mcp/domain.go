package mcp

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/url-db/internal/models"
)

type DomainManager struct {
	domainService    DomainService
	nodeCountService NodeCountService
	converter        *Converter
}

func NewDomainManager(domainService DomainService, nodeCountService NodeCountService, converter *Converter) *DomainManager {
	return &DomainManager{
		domainService:    domainService,
		nodeCountService: nodeCountService,
		converter:        converter,
	}
}

func (dm *DomainManager) ListDomains(ctx context.Context) (*MCPDomainListResponse, error) {
	response, err := dm.domainService.ListDomains(ctx, 1, 1000)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to list domains: %v", err))
	}

	mcpDomains := make([]MCPDomain, 0, len(response.Domains))
	for _, domain := range response.Domains {
		nodeCount, err := dm.nodeCountService.GetNodeCountByDomain(ctx, domain.ID)
		if err != nil {
			nodeCount = 0
		}

		mcpDomain := dm.converter.DomainToMCPDomain(&domain, nodeCount)
		if mcpDomain != nil {
			mcpDomains = append(mcpDomains, *mcpDomain)
		}
	}

	return &MCPDomainListResponse{
		Domains: mcpDomains,
	}, nil
}

func (dm *DomainManager) CreateDomain(ctx context.Context, req *models.CreateDomainRequest) (*MCPDomain, error) {
	if err := dm.validateDomainName(req.Name); err != nil {
		return nil, NewValidationError(err.Error())
	}

	normalizedName := dm.normalizeDomainName(req.Name)
	
	normalizedReq := &models.CreateDomainRequest{
		Name:        normalizedName,
		Description: req.Description,
	}

	domain, err := dm.domainService.CreateDomain(ctx, normalizedReq)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to create domain: %v", err))
	}

	return dm.converter.DomainToMCPDomain(domain, 0), nil
}

func (dm *DomainManager) GetDomain(ctx context.Context, domainName string) (*MCPDomain, error) {
	domain, err := dm.domainService.GetDomainByName(ctx, domainName)
	if err != nil {
		return nil, NewDomainNotFoundError(domainName)
	}

	nodeCount, err := dm.nodeCountService.GetNodeCountByDomain(ctx, domain.ID)
	if err != nil {
		nodeCount = 0
	}

	return dm.converter.DomainToMCPDomain(domain, nodeCount), nil
}

func (dm *DomainManager) UpdateDomain(ctx context.Context, domainName string, req *models.UpdateDomainRequest) (*MCPDomain, error) {
	domain, err := dm.domainService.GetDomainByName(ctx, domainName)
	if err != nil {
		return nil, NewDomainNotFoundError(domainName)
	}

	updatedDomain, err := dm.domainService.UpdateDomain(ctx, domain.ID, req)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to update domain: %v", err))
	}

	nodeCount, err := dm.nodeCountService.GetNodeCountByDomain(ctx, domain.ID)
	if err != nil {
		nodeCount = 0
	}

	return dm.converter.DomainToMCPDomain(updatedDomain, nodeCount), nil
}

func (dm *DomainManager) DeleteDomain(ctx context.Context, domainName string) error {
	domain, err := dm.domainService.GetDomainByName(ctx, domainName)
	if err != nil {
		return NewDomainNotFoundError(domainName)
	}

	nodeCount, err := dm.nodeCountService.GetNodeCountByDomain(ctx, domain.ID)
	if err != nil {
		return NewInternalServerError(fmt.Sprintf("failed to get node count: %v", err))
	}

	if nodeCount > 0 {
		return NewValidationError("도메인에 노드가 존재하므로 삭제할 수 없습니다")
	}

	if err := dm.domainService.DeleteDomain(ctx, domain.ID); err != nil {
		return NewInternalServerError(fmt.Sprintf("failed to delete domain: %v", err))
	}

	return nil
}

func (dm *DomainManager) GetDomainStats(ctx context.Context, domainName string) (*DomainStats, error) {
	domain, err := dm.domainService.GetDomainByName(ctx, domainName)
	if err != nil {
		return nil, NewDomainNotFoundError(domainName)
	}

	nodeCount, err := dm.nodeCountService.GetNodeCountByDomain(ctx, domain.ID)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to get node count: %v", err))
	}

	return &DomainStats{
		DomainName:  domain.Name,
		Description: domain.Description,
		NodeCount:   nodeCount,
		CreatedAt:   domain.CreatedAt,
		UpdatedAt:   domain.UpdatedAt,
	}, nil
}

func (dm *DomainManager) validateDomainName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("도메인 이름은 필수입니다")
	}

	if len(name) > 50 {
		return fmt.Errorf("도메인 이름은 50자를 초과할 수 없습니다")
	}

	matched, err := regexp.MatchString(`^[a-zA-Z0-9-]+$`, name)
	if err != nil {
		return fmt.Errorf("도메인 이름 검증 중 오류가 발생했습니다")
	}

	if !matched {
		return fmt.Errorf("도메인 이름은 영문자, 숫자, 하이픈만 사용할 수 있습니다")
	}

	if strings.HasPrefix(name, "-") || strings.HasSuffix(name, "-") {
		return fmt.Errorf("도메인 이름은 하이픈으로 시작하거나 끝날 수 없습니다")
	}

	if strings.Contains(name, "--") {
		return fmt.Errorf("도메인 이름에 연속된 하이픈은 사용할 수 없습니다")
	}

	return nil
}

func (dm *DomainManager) normalizeDomainName(name string) string {
	name = strings.ToLower(name)
	name = strings.TrimSpace(name)
	
	re := regexp.MustCompile(`[^a-z0-9-]`)
	name = re.ReplaceAllString(name, "-")
	
	re = regexp.MustCompile(`-+`)
	name = re.ReplaceAllString(name, "-")
	
	name = strings.Trim(name, "-")
	
	return name
}

func (dm *DomainManager) CheckDomainExists(ctx context.Context, domainName string) (bool, error) {
	_, err := dm.domainService.GetDomainByName(ctx, domainName)
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (dm *DomainManager) GetDomainByPartialName(ctx context.Context, partialName string) ([]MCPDomain, error) {
	response, err := dm.domainService.ListDomains(ctx, 1, 1000)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to list domains: %v", err))
	}

	var matchedDomains []MCPDomain
	partialNameLower := strings.ToLower(partialName)

	for _, domain := range response.Domains {
		if strings.Contains(strings.ToLower(domain.Name), partialNameLower) ||
			strings.Contains(strings.ToLower(domain.Description), partialNameLower) {
			
			nodeCount, err := dm.nodeCountService.GetNodeCountByDomain(ctx, domain.ID)
			if err != nil {
				nodeCount = 0
			}

			mcpDomain := dm.converter.DomainToMCPDomain(&domain, nodeCount)
			if mcpDomain != nil {
				matchedDomains = append(matchedDomains, *mcpDomain)
			}
		}
	}

	return matchedDomains, nil
}

func (dm *DomainManager) GetPopularDomains(ctx context.Context, limit int) ([]MCPDomain, error) {
	response, err := dm.domainService.ListDomains(ctx, 1, 1000)
	if err != nil {
		return nil, NewInternalServerError(fmt.Sprintf("failed to list domains: %v", err))
	}

	type domainWithCount struct {
		domain    models.Domain
		nodeCount int
	}

	var domainsWithCounts []domainWithCount

	for _, domain := range response.Domains {
		nodeCount, err := dm.nodeCountService.GetNodeCountByDomain(ctx, domain.ID)
		if err != nil {
			nodeCount = 0
		}

		domainsWithCounts = append(domainsWithCounts, domainWithCount{
			domain:    domain,
			nodeCount: nodeCount,
		})
	}

	for i := 0; i < len(domainsWithCounts)-1; i++ {
		for j := i + 1; j < len(domainsWithCounts); j++ {
			if domainsWithCounts[i].nodeCount < domainsWithCounts[j].nodeCount {
				domainsWithCounts[i], domainsWithCounts[j] = domainsWithCounts[j], domainsWithCounts[i]
			}
		}
	}

	if limit > 0 && limit < len(domainsWithCounts) {
		domainsWithCounts = domainsWithCounts[:limit]
	}

	var popularDomains []MCPDomain
	for _, dwc := range domainsWithCounts {
		mcpDomain := dm.converter.DomainToMCPDomain(&dwc.domain, dwc.nodeCount)
		if mcpDomain != nil {
			popularDomains = append(popularDomains, *mcpDomain)
		}
	}

	return popularDomains, nil
}

type DomainStats struct {
	DomainName  string    `json:"domain_name"`
	Description string    `json:"description"`
	NodeCount   int       `json:"node_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

import "time"