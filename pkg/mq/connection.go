package mq

import (
	"strconv"

	"github.com/isayme/go-amqp-reconnect/rabbitmq"
)

const (
	timeout = "connection_timeout"
)

func NewConnectionWithConfig(conf *Config) (*rabbitmq.Connection, error) {
	t := conf.Data().Query().Get(timeout)
	if t == "" {
		conf.Data().Query().Set(timeout, strconv.Itoa(defaultConnectTimeout))
	}
	conn, err := rabbitmq.Dial(conf.Dsn())
	if err != nil {
		return nil, err
	}

	return conn, nil
}
