package extensions

import (
	"log"
	dataContext "urlshortener-crud-consumer/data"
	"urlshortener-crud-consumer/data/repositories"
	"urlshortener-crud-consumer/services"
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

	elasticSearchService, err := service.NewElasticSearchService("", []string{""}, 3)
	if err != nil {
		log.Panic(err)
	}

	queueService := service.NewQueueService()

	urlShorthenerService := service.NewUrlShorthenerService(repo, elasticSearchService, queueService)

	return &Services{UrlShorthener: urlShorthenerService}
}
