package mq

import (
	"strconv"

	"github.com/isayme/go-amqp-reconnect/rabbitmq"
)

func NewConnectionWithConfig(conf Config) (*rabbitmq.Connection, error) {
	if conf.Param(paramConnectionTimeout) == "" {
		conf.Params[paramConnectionTimeout] = strconv.Itoa(defaultConnectTimeout)
	}
	conn, err := rabbitmq.Dial(conf.Dsn())
	if err != nil {
		return nil, err
	}

	return conn, nil
}
