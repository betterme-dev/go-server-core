package healtcheck

import (
	"errors"

	"github.com/heptiolabs/healthcheck"

	"github.com/betterme-dev/go-server-core/pkg/env"
)

func RabbitMQCheck() healthcheck.Check {
	return func() error {
		mqConn := env.Queue()
		if mqConn == nil {
			return errors.New("RabbitMQ connection is nil")
		}
		if mqConn.IsConnectionClosed() {
			return errors.New("RabbitMQ connection closed")
		}

		return nil
	}
}

func Neo4jCheck() healthcheck.Check {
	return func() error {
		panic("Implement me")
	}
}
