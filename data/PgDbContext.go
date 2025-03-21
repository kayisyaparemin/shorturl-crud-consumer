package dataContext

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func ConnectDB() (*gorm.DB,error) {
	dsn := ""
	err := GetJSON("url",&dsn)
	if err != nil {
		return nil,err
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil,err
	}
	return db,nil
}

func GetJSON[T any](url string, result *T) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("request Fail: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error Code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("response was not readed: %w", err)
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		return fmt.Errorf("json parse error: %w", err)
	}

	return nil
}
