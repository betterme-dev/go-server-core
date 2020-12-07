package healtcheck

import (
	"time"

	hatcheckLib "github.com/heptiolabs/healthcheck"
	"github.com/spf13/viper"

	"github.com/betterme-dev/go-server-core/pkg/env"
)

type Checks []hatcheckLib.Check

const DBTimeout = 1 * time.Second
const ESTimeout = 1 * time.Second

/**
Usage example:
	checks := healtcheck.Checks{
		healtcheck.DB(),
		healtcheck.ElasticSearch(),
		healtcheck.RabbitMQ(),
	}
	probs := healtcheck.ConfigHandler(checks, checks)

  With our Gin app:
	app := web.NewApp()
	app.Engine.Handle("GET", "/ready", gin.WrapF(probs.ReadyEndpoint))
	app.Engine.Handle("GET", "/live", gin.WrapF(probs.LiveEndpoint))

*/
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
	return hatcheckLib.DatabasePingCheck(env.DB(), DBTimeout)
}

func ElasticSearch() hatcheckLib.Check {
	return hatcheckLib.HTTPGetCheck(viper.GetString("ELASTICSEARCH_ADDRESS"), ESTimeout)
}

func RabbitMQ() hatcheckLib.Check {
	return rabbitMQCheck()
}

func Neo4j() hatcheckLib.Check {
	return neo4jCheck()
}

func Redis() hatcheckLib.Check {
	return redisCheck()
}
