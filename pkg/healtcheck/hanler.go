package healtcheck

import (
	"time"

	hatcheckLib "github.com/heptiolabs/healthcheck"
	"github.com/spf13/viper"

	"github.com/betterme-dev/go-server-core/pkg/env"
)

type Checks []hatcheckLib.Check

func ConfigHandler(readiness Checks, liveness Checks) hatcheckLib.Handler {
	handler := hatcheckLib.NewHandler()
	for name, check := range readiness {
		handler.AddReadinessCheck(string(rune(name)), check)
	}

	for name, check := range liveness {
		handler.AddLivenessCheck(string(rune(name)), check)
	}

	return handler
}

func DB() hatcheckLib.Check {
	return hatcheckLib.DatabasePingCheck(env.DB(), 1*time.Second)
}

func ElasticSearch() hatcheckLib.Check {
	return hatcheckLib.HTTPGetCheck(viper.GetString("ELASTICSEARCH_ADDRESS"), 1*time.Second)
}

func RabbitMQ() hatcheckLib.Check {
	return RabbitMQCheck()
}

func Neo4j() hatcheckLib.Check {
	return Neo4jCheck()
}
