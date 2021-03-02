package mq

import (
	"errors"
	"net/url"
)

type Config struct {
	data *url.URL
}

func NewConfig(dsn string) (*Config, error) {
	conf, err := parseDsn(dsn)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func (c Config) Dsn() string {
	if c.data == nil {
		return ""
	}
	return c.Data().String()
}

func (c Config) Data() *url.URL {
	return c.data
}
func parseDsn(dsn string) (*Config, error) {
	if dsn == "" {
		return nil, errors.New("dsn string is empty")
	}
	var conf Config
	data, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}
	conf.data = data
	return &conf, nil
}
