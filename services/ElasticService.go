package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"github.com/elastic/go-elasticsearch/v8"
	response "urlshortener-crud-consumer/models/responses"
)
type IElasticSearchService interface{
	CreateIndex() (string, error)
	BulkAddNewItemsToIndex(data []response.ShortUrlSearch, indexName string) error
	GetIndexName() (string, error) 
}

type ElasticSearchService struct {
	client    *elasticsearch.Client
	aliasName string
}

func NewElasticSearchService(env string, urls []string, timeout time.Duration) (*ElasticSearchService, error) {
	cfg := elasticsearch.Config{
		Addresses: urls,
		Transport: nil,
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	aliasName := fmt.Sprintf("shorturl-%s", strings.ToLower(env))
	return &ElasticSearchService{
		client:    client,
		aliasName: aliasName,
	}, nil
}

func (es *ElasticSearchService) CreateIndex() (string, error) {
	indexName := fmt.Sprintf("%s-%d", es.aliasName, time.Now().UnixNano())

	res, err := es.client.Indices.Exists([]string{indexName})
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		return "", fmt.Errorf("index %s already exists", indexName)
	}

	settings := map[string]interface{}{
		"settings": map[string]interface{}{
			"number_of_shards":   1,
			"number_of_replicas": 1,
		},
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"key":             map[string]string{"type": "keyword"},
				"long_url":        map[string]string{"type": "text"},
				"telephoneNumber": map[string]string{"type": "keyword"},
			},
		},
	}

	body, err := json.Marshal(settings)
	if err != nil {
		return "", err
	}

	res, err = es.client.Indices.Create(indexName, es.client.Indices.Create.WithBody(bytes.NewReader(body)))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.IsError() {
		return "", fmt.Errorf("error creating index: %s", res.String())
	}

	aliasReq := map[string]interface{}{
		"actions": []map[string]interface{}{
			{"add": map[string]string{"index": indexName, "alias": es.aliasName}},
		},
	}

	body, err = json.Marshal(aliasReq)
	if err != nil {
		return "", err
	}

	res, err = es.client.Indices.UpdateAliases(bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.IsError() {
		return "", fmt.Errorf("error updating alias: %s", res.String())
	}

	return indexName, nil
}

func (es *ElasticSearchService) BulkAddNewItemsToIndex(data []response.ShortUrlSearch, indexName string) error {
	if len(data) == 0 {
		return nil 
	}

	var buf bytes.Buffer

	for _, item := range data {
		meta := map[string]interface{}{
			"index": map[string]string{"_index": indexName},
		}

		metaJSON, err := json.Marshal(meta)
		if err != nil {
			return err
		}

		itemJSON, err := json.Marshal(item)
		if err != nil {
			return err
		}

		buf.Write(metaJSON)
		buf.WriteString("\n")
		buf.Write(itemJSON)
		buf.WriteString("\n")
	}
	res, err := es.client.Bulk(bytes.NewReader(buf.Bytes()), es.client.Bulk.WithContext(context.Background()))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("bulk insert error: %s", res.String())
	}

	fmt.Println("Success: Documents were added to Index.")
	return nil
}

func (es *ElasticSearchService) GetIndexName() (string, error) {
	res, err := es.client.Indices.GetAlias(es.client.Indices.GetAlias.WithContext(context.Background()))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.IsError() {
		return "", fmt.Errorf("error getting alias: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return "", err
	}

	for index := range result {
		if strings.HasPrefix(index, es.aliasName) {
			return index, nil
		}
	}
	return "", fmt.Errorf("no index found for alias %s", es.aliasName)
}
