package mq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConnectionWithConfig(t *testing.T) {
	// prepare
	conf := NewConfig()
	conf.Username = "rabbitadmin"
	conf.Host = "rabbitmq-eks.default.svc.cluster.local"
	conf.Vhost = "/better-me-v3"
	conf.Params["connection_timeout"] = "2000"

	// raw password
	conf.Password = "gL4]}J@A^juSxD]x"
	conn, err := NewConnectionWithConfig(conf)
	assert.Nil(t, err)
	assert.NotNil(t, conn)
	// encoded password
	conf.Password = "gL4%5D%7DJ%40A%5EjuSxD%5Dx"
	conn, err = NewConnectionWithConfig(conf)
	assert.Nil(t, err)
	assert.NotNil(t, conn)
}
