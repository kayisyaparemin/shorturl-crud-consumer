package utilities

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

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
