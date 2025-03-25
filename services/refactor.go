package service

import (
	"log"
	"sync"
	"urlshortener-crud-consumer/data/repositories"
	"urlshortener-crud-consumer/models"
	"urlshortener-crud-consumer/models/entities"
	"urlshortener-crud-consumer/models/requests"
	"urlshortener-crud-consumer/models/responses"
	"urlshortener-crud-consumer/utilities"
)

// Yeni constructor fonksiyonu
func NewUrlShorthenerService2(repo repositories.IUrlShortenerRepository, elastic IElasticSearchService, queue IQueueService) *UrlShorthenerService {
	return &UrlShorthenerService{repo: repo, Elastic: elastic, Queue: queue}
}

func (uss *UrlShorthenerService) UpdateShortUrls(model *requests.QueueModel) {
	const numWorkers = 1000
	const partialSize = 200

	taskCount := (len(model.TelephoneNumbers) + partialSize - 1) / partialSize

	indexName, err := uss.Elastic.CreateIndex()
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	chn := make(chan models.UpdateUrlShortenerModel, numWorkers)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			uss.Update(chn, indexName, model.ChannelCampaignId, model.PartialId, model.Url, model.UserId)
		}()
	}
	for i := 0; i < taskCount; i++ {
		keysPartition, err := uss.repo.GetAvailableKeys(partialSize)
		if err != nil {
			log.Printf("Error: %v\n", err)
		}
		start := i * partialSize
		end := start + partialSize
		if end > len(model.TelephoneNumbers) {
			end = len(model.TelephoneNumbers)
		}
		phonePartition := model.TelephoneNumbers[start:end]
		updateModel := models.UpdateUrlShortenerModel{
			Phones: phonePartition,
			Keys:   keysPartition,
		}

		chn <- updateModel
	}

	close(chn)

	wg.Wait()
}

func (uss *UrlShorthenerService) Update(chn <-chan models.UpdateUrlShortenerModel, indexName string, channelCampaignId int, partialId int, longUrl string, userId int) {
	for updateModel := range chn { 
		var bulkUpdateModels []entities.UrlShortener
		var bulkElasticUpdates []responses.ShortUrlSearch

		for j := 0; j < len(updateModel.Phones) && j < len(updateModel.Keys); j++ {
			urlShortener := updateModel.Keys[j]
			phone := updateModel.Phones[j]

			bulkUpdateModels = append(bulkUpdateModels, utilities.UpdateUrlShortenerModel(urlShortener, channelCampaignId, partialId, longUrl, phone, userId))
			bulkElasticUpdates = append(bulkElasticUpdates, responses.ShortUrlSearch{Key: urlShortener.Key, LongUrl: longUrl, TelephoneNumber: phone})
		}
		var subWg sync.WaitGroup
		subWg.Add(2)

		go func() {
			defer subWg.Done()
			if err := uss.Elastic.BulkAddNewItemsToIndex(bulkElasticUpdates, indexName); err != nil {
				log.Printf("Error updating ElasticSearch: %v", err)
			}
		}()

		go func() {
			defer subWg.Done()
			if err := uss.repo.BulkUpdateAvailableToUsedKeys(bulkUpdateModels); err != nil {
				log.Printf("Error updating database: %v", err)
			}
		}()

		subWg.Wait()
	}
}
