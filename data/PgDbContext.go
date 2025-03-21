package dataContext

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"urlshortener-crud-consumer/models/responses"
    "strings"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	response := new(responses.ConfResponse)
	err := GetJSON("http://gatewayapi.test-gateways/configuration-services/configurations/hangikredi.shorturlservice.postgres.connectionstring", response)
	if err != nil {
		return nil, err
	}
	dsn := GetDsn(response.Value)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetJSON(url string, result *responses.ConfResponse) error {
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
func GetDsn(dsn string) (string){
params := strings.Split(dsn, ";")
	
connParams := make(map[string]string)

for _, param := range params {
	parts := strings.Split(param, "=")
	if len(parts) == 2 {
		key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		connParams[key] = value
	}
}
connectionString := fmt.Sprintf("host=%s user=%s password='%s' dbname=%s port=%s sslmode=disable TimeZone=UTC",
	connParams["HOST"],
	connParams["USER ID"],
	connParams["PASSWORD"],
	connParams["DATABASE"],
	connParams["PORT"])
return connectionString
}
