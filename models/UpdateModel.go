package models

import "urlshortener-crud-consumer/models/entities"

type UpdateUrlShortenerModel struct {
	Phones []string 
	Keys []entities.UrlShortener
}