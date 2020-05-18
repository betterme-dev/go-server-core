package elasticsearch

import (
	"fmt"
	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Client() (*elasticsearch7.Client, error) {
	cfg := elasticsearch7.Config{
		Addresses: []string{
			viper.GetString("ELASTICSEARCH_ADDRESS"),
		},
	}

	es, err := elasticsearch7.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("can't create elasticsearch client: %s", err)
	}

	log.Info("Elasticsearch client created ")

	return es, nil
}
