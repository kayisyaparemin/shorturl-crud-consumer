package tests

import (
	"sync"
	"testing"
	"urlshortener-crud-consumer/models/entities"
	service "urlshortener-crud-consumer/services"

	"github.com/stretchr/testify/mock"
)

// Test
func TestUpdateShortUrl(t *testing.T) {
	// Mock objects
	mockRepo := new(MockUrlShortenerRepository)
	mockElastic := new(MockElasticSearchService)

	// Test data
	telephones := []string{"12345", "67890"}
	keys := []entities.UrlShortener{
		{Key: "key1", Id: 1},
		{Key: "key2", Id: 2},
	}
	indexName := "test_index"
	channelCampaignId := 101
	partialId := 202
	longUrl := "http://example.com"
	userId := 1001

	// Expectations for mock objects
	mockRepo.On("BulkUpdateAvailableToUsedKeys", mock.Anything).Return(nil).Once()
	mockElastic.On("BulkAddNewItemsToIndex", mock.Anything, indexName).Return(nil).Once()

	// Create service with mock objects
	service := service.NewUrlShorthenerService(mockRepo, mockElastic, nil)

	// WaitGroup for concurrency
	var wg sync.WaitGroup

	// Add goroutine count to WaitGroup (1 in this case)
	wg.Add(1)

	// Call the method to test in a goroutine
	go service.UpdateShortUrl(telephones, keys, indexName, channelCampaignId, partialId, longUrl, userId, &wg)

	// Wait for all goroutines to finish
	wg.Wait()

	// Assertions for mocks
	mockRepo.AssertExpectations(t)
	mockElastic.AssertExpectations(t)
}
