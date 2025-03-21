package tests

import (
	"sync"
	"testing"
	"urlshortener-crud-consumer/models/entities"
	service "urlshortener-crud-consumer/services"
	"github.com/stretchr/testify/mock"
)

func TestCreateShortUrl(t *testing.T) {
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
	partialSize := 1
	taskCount := 2

	// Expectations for mock objects
	mockRepo.On("BulkUpdateAvailableToUsedKeys", mock.Anything).Return(nil).Times(taskCount) 
	mockElastic.On("BulkAddNewItemsToIndex", mock.Anything, indexName).Return(nil).Times(taskCount)

	// Create service with mock objects
	service := service.NewUrlShorthenerService(mockRepo, mockElastic, nil)

	// Channels for partitioning
	telephoneChannel := make(chan []string, taskCount)
	keyChannel := make(chan []entities.UrlShortener, taskCount)

	// Simulate partitioning
	go service.StreamingPartitionUrlShortener(keys, partialSize, keyChannel)
	go service.StreamingPartitionString(telephones, partialSize, telephoneChannel)

	// WaitGroup for concurrency
	var wg sync.WaitGroup
	for i := 0; i < taskCount; i++ {
		telephones, ok1 := <-telephoneChannel
		keys, ok2 := <-keyChannel
		if !ok1 || !ok2 {
			break
		}
		wg.Add(1)
		go service.UpdateShortUrl(telephones, keys, indexName, channelCampaignId, partialId, longUrl, userId, &wg)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Assertions
	mockRepo.AssertExpectations(t)   
	mockElastic.AssertExpectations(t)
}

