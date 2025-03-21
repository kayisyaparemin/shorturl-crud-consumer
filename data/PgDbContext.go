package dataContext

import (
	"urlshortener-crud-consumer/models/responses"
	"urlshortener-crud-consumer/utilities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	response := new(responses.DbSettings)
	err := utilities.GetJSON("http://gatewayapi.test-gateways/configuration-services/configurations/hangikredi.shorturlservice.postgres.connectionstring", response)
	if err != nil {
		return nil, err
	}
	dsn := utilities.GetDsn(response.Value)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}


