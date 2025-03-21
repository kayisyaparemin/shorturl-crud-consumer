package main

import (
	"urlshortener-crud-consumer/consumers"
	"urlshortener-crud-consumer/core/extensions"
)
func main(){
	services := extensions.RegisterServices()
	consumer.Consume(services)
}