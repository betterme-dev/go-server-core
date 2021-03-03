package mq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	conf := NewConfig()
	assert.NotNil(t, conf)
	assert.Equal(t, defaultDriver, conf.Driver)
	assert.Equal(t, defaultHost, conf.Host)
	assert.Equal(t, defaultPort, conf.Port)
}

func TestNewConfigFromDsnSuccess(t *testing.T) {
	dsn := "amqp://user:pass@host.com:1234//vhost?param=12"
	conf, err := NewConfigFromDsn(dsn)
	// struct check
	assert.Nil(t, err)
	assert.NotNil(t, conf)
	assert.IsType(t, Config{}, conf)
	// data check
	assert.Equal(t, "amqp", conf.Driver)
	assert.Equal(t, "user", conf.Username)
	assert.Equal(t, "user", conf.UsernameEscaped())
	assert.Equal(t, "pass", conf.Password)
	assert.Equal(t, "pass", conf.PasswordEscaped())
	assert.Equal(t, "host.com", conf.Host)
	assert.Equal(t, "1234", conf.Port)
	assert.Equal(t, "/vhost", conf.Vhost)
	assert.Equal(t, "%2Fvhost", conf.VhostEscaped())
	assert.Equal(t, "12", conf.Param("param"))
	// dsn check
	assert.Equal(t, dsn, conf.Dsn())
	conf.Params["timeout"] = "3"
	assert.Equal(t, "3", conf.Param("timeout"))
	// params order is unpredicted, so...
	res1 := conf.Dsn() == "amqp://user:pass@host.com:1234//vhost?param=12&timeout=3"
	res2 := conf.Dsn() == "amqp://user:pass@host.com:1234//vhost?timeout=3&param=12"
	assert.True(t, res1 || res2)
}

func TestNewConfigFromDsnFail(t *testing.T) {
	// scenario: empty dsn
	conf, err := NewConfigFromDsn("")
	assert.NotNil(t, err)
	// scenario: bad dsn
	conf, err = NewConfigFromDsn("amqp://user:@host:1234//vhost")
	assert.Nil(t, err)
	assert.NotNil(t, conf)
	assert.IsType(t, Config{}, conf)
	assert.Equal(t, "user", conf.Username)
	assert.Equal(t, "", conf.Password)
	assert.Equal(t, "host", conf.Host)
	assert.Equal(t, "1234", conf.Port)
	assert.NotEqual(t, "vhost", conf.Vhost)
	assert.NotEqual(t, "vhost", conf.VhostEscaped())
	assert.Equal(t, "", conf.Param("param"))

	// scenario: unescaped auth data
	dsn := "amqp://us$^er:fdjkh%23${#!fdf0-0&@rabbitmq.host:5672//my-vhost?connection_timeout=500"
	conf, err = NewConfigFromDsn(dsn)
	assert.NotNil(t, err)
	dsnOK := "amqp://us%24%5Eer:fdjkh%2523%24%7B%23%21fdf0-0%26@rabbitmq.host:5672//my-vhost?connection_timeout=500"
	conf, err = NewConfigFromDsn(dsnOK)
	assert.Nil(t, err)
}
