package service

import (
	"encoding/json"
	"log"
	"sync"
	"urlshortener-crud-consumer/data/repositories"
	"urlshortener-crud-consumer/models/entities"
	"urlshortener-crud-consumer/models/queues"
	"urlshortener-crud-consumer/models/requests"
	"urlshortener-crud-consumer/models/responses"
	"urlshortener-crud-consumer/utilities"
)

type UrlShorthenerService struct {
	repo    repositories.IUrlShortenerRepository
	Elastic IElasticSearchService
	Queue   IQueueService
}

// Yeni constructor fonksiyonu, interface alacak şekilde güncellendi
func NewUrlShorthenerService(repo repositories.IUrlShortenerRepository, elastic IElasticSearchService, queue IQueueService) *UrlShorthenerService {
	return &UrlShorthenerService{repo: repo, Elastic: elastic, Queue: queue}
}

func (uss *UrlShorthenerService) StartProcess(model *requests.QueueModel) {
	go uss.createShortUrlAsync(model)
}

func (uss *UrlShorthenerService) createShortUrlAsync(model *requests.QueueModel) {
	if model.Email == "" {
		model.Email = "Teknoloji-Transformers@hangikredi.com"
	}
	extraData := map[string]interface{}{"eventTriggerId": 8598}
	extraParameterJSON, _ := json.Marshal(extraData)
	uss.Queue.Send(queues.MailQueueModel{
		To:              []string{model.Email},
		ServiceProvider: queues.ServiceProviderEnumType(queues.Emarsys),
		ExtraParameters: extraParameterJSON,
		Subject:         "Kampanya Excel Yüklemesi",
		Body:            "Excel yüklenmesine başlandı.",
	})

	const partialSize = 2000

	keysPartition, err := uss.repo.GetAvailableKeys(model.Size)
	if err != nil {
		log.Panic(err)
	}

	indexName, err := uss.Elastic.CreateIndex()
	if err != nil {
		log.Panic(err)
	}
	taskCount := (len(model.TelephoneNumbers) + partialSize - 1) / partialSize

	telephoneChannel := make(chan []string, taskCount)
	keyChannel := make(chan []entities.UrlShortener, taskCount)

	go uss.StreamingPartitionUrlShortener(keysPartition, partialSize, keyChannel)
	go uss.StreamingPartitionString(model.TelephoneNumbers, partialSize, telephoneChannel)

	var wg sync.WaitGroup
	for i := 0; i < taskCount; i++ {
		telephones, ok1 := <-telephoneChannel
		keys, ok2 := <-keyChannel
		if !ok1 || !ok2 {
			break
		}
		wg.Add(1)
		go uss.UpdateShortUrl(telephones, keys, indexName, model.ChannelCampaignId, model.PartialId, model.Url, model.UserId, &wg)
	}

	go func() {
		wg.Wait()
		close(telephoneChannel)
		close(keyChannel)
	}()
	uss.Queue.Send(queues.MailQueueModel{
		To:              []string{model.Email},
		ServiceProvider: queues.ServiceProviderEnumType(queues.Emarsys),
		ExtraParameters: extraParameterJSON,
		Subject:         "Kampanya Excel Yüklemesi",
		Body:            "Excel yükleme tamamlandı.",
	})
}

func (uss *UrlShorthenerService) UpdateShortUrl(telephones []string, keys []entities.UrlShortener, indexName string, channelCampaignId int, partialId int, longUrl string, userId int, wg *sync.WaitGroup) {
	defer wg.Done()

	var bulkUpdateModels []entities.UrlShortener
	var bulkElasticUpdates []responses.ShortUrlSearch

	for j := 0; j < len(telephones) && j < len(keys); j++ {
		urlShortener := keys[j]
		phone := telephones[j]

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

func (uss *UrlShorthenerService) StreamingPartitionString(input []string, batchSize int, out chan<- []string) {
	defer close(out)
	for i := 0; i < len(input); i += batchSize {
		end := i + batchSize
		if end > len(input) {
			end = len(input)
		}
		out <- input[i:end]
	}
}

func (uss *UrlShorthenerService) StreamingPartitionUrlShortener(input []entities.UrlShortener, batchSize int, out chan<- []entities.UrlShortener) {
	defer close(out)
	for i := 0; i < len(input); i += batchSize {
		end := i + batchSize
		if end > len(input) {
			end = len(input)
		}
		out <- input[i:end]
	}
}
