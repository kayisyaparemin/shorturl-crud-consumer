package responses

import (
	"strings"
	"strconv"
)

type ElasticSearchSettings struct {
	NumberOfReplicas int    `json:"numberofreplicas"`
	NumberOfShards   int    `json:"numberofshards"`
	TimeOut          int    `json:"timeout"`
	Url              string `json:"url"`
}

func GetElasticSearchSettings(s *[]Settings) *ElasticSearchSettings {

	settings := new(ElasticSearchSettings)
	for _, setting := range *s{

		key := strings.Split(setting.Key, ".")[1]
		switch key {
		case "url":
			settings.Url = setting.Value
		case "timeout":
			settings.TimeOut,_ = strconv.Atoi(setting.Value)
		}
	}
	return settings
}
