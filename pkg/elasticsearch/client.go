package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ES struct {
	Client *elasticsearch7.Client
}

func NewClient() (es *ES, err error) {
	cfg := elasticsearch7.Config{
		Addresses: []string{
			viper.GetString("ELASTICSEARCH_ADDRESS"),
		},
	}

	es = &ES{}
	es.Client, err = elasticsearch7.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("can't create elasticsearch client: %s", err)
	}

	log.Info("Elasticsearch client created ")

	return
}

func (es *ES) DoSearch(query map[string]interface{}, index string) (map[string]interface{}, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("error encoding query: %s", err)
	}

	client := es.Client
	res, err := client.Search(
		client.Search.WithContext(context.Background()),
		client.Search.WithIndex(index),
		client.Search.WithBody(&buf),
		client.Search.WithTrackTotalHits(true),
		client.Search.WithPretty(),
	)
	if err != nil {
		return nil, fmt.Errorf("error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("error parsing the response body: %s", err)
		}
		// Print the response status and error information.
		return nil, fmt.Errorf("[%s] %s: %s",
			res.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"],
		)
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("error parsing the response body: %s", err)
	}

	return r, err
}
