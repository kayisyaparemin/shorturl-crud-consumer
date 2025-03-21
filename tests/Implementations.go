package tests

import (
	"urlshortener-crud-consumer/models/entities"
	"urlshortener-crud-consumer/models/responses"

	"github.com/stretchr/testify/mock"
)

// Mock Repository
type MockUrlShortenerRepository struct {
	mock.Mock
}

func (m *MockUrlShortenerRepository) BulkUpdateAvailableToUsedKeys(entityModels []entities.UrlShortener) error {
	args := m.Called(entityModels)
	return args.Error(0)
}

func (m *MockUrlShortenerRepository) GetAvailableKeys(size int) ([]entities.UrlShortener, error) {
	args := m.Called(size)
	return args.Get(0).([]entities.UrlShortener), args.Error(1)
}

// Mock ElasticSearchService
type MockElasticSearchService struct {
	mock.Mock
}

func (m *MockElasticSearchService) BulkAddNewItemsToIndex(bulkElasticUpdates []responses.ShortUrlSearch, indexName string) error {
	args := m.Called(bulkElasticUpdates, indexName)
	return args.Error(0)
}

func (m *MockElasticSearchService) CreateIndex() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockElasticSearchService) GetIndexName() (string, error) {
	args := m.Called()
	return args.String(0), nil
}
