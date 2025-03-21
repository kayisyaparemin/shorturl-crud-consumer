package dataContext

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB,error) {
	dsn := "host=localhost user=admin password=password dbname=userdb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil,err
	}
	return db,nil
}
