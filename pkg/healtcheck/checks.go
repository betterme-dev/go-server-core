package healtcheck

import (
	"fmt"

	"github.com/heptiolabs/healthcheck"

	"github.com/betterme-dev/go-server-core/pkg/env"
	"github.com/betterme-dev/go-server-core/pkg/neo4j"
)

func rabbitMQCheck() healthcheck.Check {
	return func() error {
		mqConn := env.Queue()
		if mqConn == nil {
			return fmt.Errorf("rabbitMQ connection is nil")
		}
		if mqConn.IsConnectionClosed() {
			return fmt.Errorf("rabbitMQ connection closed")
		}

		return nil
	}
}

func neo4jCheck() healthcheck.Check {
	return func() error {
		c, err := neo4j.NewClient()
		if c == nil || err != nil {
			return fmt.Errorf("can`t connect to nei4j")
		}

		return nil
	}
}

func redisCheck() healthcheck.Check {
	return func() error {
		panic("implement me")
	}
}
