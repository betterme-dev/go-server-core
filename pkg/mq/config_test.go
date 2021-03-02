package mq

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigSuccess(t *testing.T) {
	dsn := "amqp://user:pass@host.com:1234/vhost?param=12"
	conf, err := NewConfig(dsn)
	// struct check
	assert.Nil(t, err)
	assert.NotNil(t, conf)
	assert.IsType(t, &Config{}, conf)
	assert.NotNil(t, conf.Data())
	// data check
	assert.IsType(t, &url.URL{}, conf.Data())
	assert.Equal(t, "amqp", conf.Data().Scheme)
	assert.Equal(t, "user", conf.Data().User.Username())
	pass, set := conf.Data().User.Password()
	assert.True(t, set)
	assert.Equal(t, "pass", pass)
	assert.Equal(t, "host.com", conf.Data().Hostname())
	assert.Equal(t, "1234", conf.Data().Port())
	assert.Equal(t, "/vhost", conf.Data().Path)
	assert.Equal(t, "12", conf.Data().Query().Get("param"))
	// dsn check
	assert.Equal(t, dsn, conf.Data().String())
	q := conf.Data().Query()
	q.Set("timeout", "3")
	conf.Data().RawQuery = q.Encode()
	assert.Equal(t, "3", conf.Data().Query().Get("timeout"))
	assert.Equal(t, conf.Dsn(), "amqp://user:pass@host.com:1234/vhost?param=12&timeout=3")
}

func TestNewConfigFail(t *testing.T) {
	// scenario: empty dsn
	conf, err := NewConfig("")
	assert.NotNil(t, err)
	assert.Nil(t, conf)

	// scenario: bad dsn
	conf, err = NewConfig("amqp://user:@host:1234//vhost")
	assert.Nil(t, err)
	assert.NotNil(t, conf)
	assert.IsType(t, &Config{}, conf)
	assert.NotNil(t, conf.Data())
	assert.IsType(t, &url.URL{}, conf.Data())
	assert.Equal(t, "amqp", conf.Data().Scheme)
	assert.Equal(t, "user:", conf.Data().User.String())
	assert.Equal(t, "user", conf.Data().User.Username())
	pass, set := conf.Data().User.Password()
	assert.True(t, set)
	assert.Equal(t, "", pass)
	assert.Equal(t, "host", conf.Data().Hostname())
	assert.Equal(t, "1234", conf.Data().Port())
	assert.NotEqual(t, "/vhost", conf.Data().Path)
	assert.Equal(t, "", conf.Data().Query().Get("param"))

	// scenario: valid, but unescaped dsn
	dsn := "amqp://user:fdjkh%23${#!fdf0-0&@rabbitmq.host:5672//my-vhost?connection_timeout=500"
	conf, err = NewConfig(dsn)
	assert.NotNil(t, err)
	assert.Nil(t, conf)
	escapedDSN := url.PathEscape(dsn)
	conf, err = NewConfig(escapedDSN)
	assert.Nil(t, err)
	assert.NotNil(t, conf)
}
