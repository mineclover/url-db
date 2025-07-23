package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"url-db/internal/models"
)

func TestDomainService_CreateDomain(t *testing.T) {
	service, mockRepo := CreateTestDomainService(t)

	ctx := CreateTestContext()
	req := &models.CreateDomainRequest{
		Name:        "test-domain",
		Description: "Test domain description",
	}

	mockRepo.SetShouldError(false)

	domain, err := service.CreateDomain(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, domain)
	assert.Equal(t, "test-domain", domain.Name)
	assert.Equal(t, "Test domain description", domain.Description)
}

func TestDomainService_GetDomain(t *testing.T) {
	service, mockRepo := CreateTestDomainService(t)

	ctx := CreateTestContext()
	testDomain := CreateTestDomain("test", "Test description")
	mockRepo.Create(ctx, testDomain)

	domain, err := service.GetDomain(ctx, testDomain.ID)

	assert.NoError(t, err)
	assert.NotNil(t, domain)
	assert.Equal(t, testDomain.ID, domain.ID)
	assert.Equal(t, testDomain.Name, domain.Name)
}

func TestDomainService_ListDomains(t *testing.T) {
	service, mockRepo := CreateTestDomainService(t)

	ctx := CreateTestContext()
	domain1 := CreateTestDomain("domain1", "Description 1")
	domain2 := CreateTestDomain("domain2", "Description 2")

	mockRepo.Create(ctx, domain1)
	mockRepo.Create(ctx, domain2)

	response, err := service.ListDomains(ctx, 1, 20)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 2, response.TotalCount)
	assert.Len(t, response.Domains, 2)
}

func TestDomainService_UpdateDomain(t *testing.T) {
	service, mockRepo := CreateTestDomainService(t)

	ctx := CreateTestContext()
	testDomain := CreateTestDomain("test", "Original description")
	mockRepo.Create(ctx, testDomain)

	req := &models.UpdateDomainRequest{
		Description: "Updated description",
	}

	domain, err := service.UpdateDomain(ctx, testDomain.ID, req)

	assert.NoError(t, err)
	assert.NotNil(t, domain)
	assert.Equal(t, "Updated description", domain.Description)
}

func TestDomainService_DeleteDomain(t *testing.T) {
	service, mockRepo := CreateTestDomainService(t)

	ctx := CreateTestContext()
	testDomain := CreateTestDomain("test", "Test description")
	mockRepo.Create(ctx, testDomain)

	err := service.DeleteDomain(ctx, testDomain.ID)

	assert.NoError(t, err)

	// Verify it's deleted
	_, err = service.GetDomain(ctx, testDomain.ID)
	assert.Error(t, err)
}
