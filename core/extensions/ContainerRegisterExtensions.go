package extensions

import (
	"log"
	dataContext "urlshortener-crud-consumer/data"
	"urlshortener-crud-consumer/data/repositories"
	"urlshortener-crud-consumer/models/responses"
	"urlshortener-crud-consumer/services"
	"urlshortener-crud-consumer/utilities"
	"strings"
)

type Services struct {
	UrlShorthener *service.UrlShorthenerService
}

func RegisterServices() *Services {
	database, err := dataContext.ConnectDB()
	if err != nil {
		log.Panic(err)
	}

	repo := repositories.NewUrlShortenerRepository(database)

	response := new([]responses.Settings)
	utilities.GetJSON("http://configurationapi.test-microservices/api/configurations/elasticsearchsettings/list",response)
	elasticSettings := responses.GetElasticSearchSettings(response)

	elasticSearchService, err := service.NewElasticSearchService("test", strings.Split(elasticSettings.Url,","), elasticSettings.TimeOut)
	if err != nil {
		log.Panic(err)
	}

	queueService := service.NewQueueService()

	urlShorthenerService := service.NewUrlShorthenerService(repo, elasticSearchService, queueService)

	return &Services{UrlShorthener: urlShorthenerService}
}
