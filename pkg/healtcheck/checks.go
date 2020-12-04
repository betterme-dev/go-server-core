package healtcheck

import (
	"errors"

	"github.com/heptiolabs/healthcheck"

	"github.com/betterme-dev/go-server-core/pkg/env"
	"github.com/betterme-dev/go-server-core/pkg/neo4j"
)

func RabbitMQCheck() healthcheck.Check {
	return func() error {
		mqConn := env.Queue()
		if mqConn == nil {
			return errors.New("rabbitMQ connection is nil")
		}
		if mqConn.IsConnectionClosed() {
			return errors.New("rabbitMQ connection closed")
		}
		return nil
	}
}

func Neo4jCheck() healthcheck.Check {
	return func() error {
		c, err := neo4j.NewClient()
		if c == nil || err != nil {
			return errors.New("can`t connect to nei4j")
		}
		return nil
	}
}
