package mq

import "net/url"

type Config struct {
	data *url.URL
}

func NewConfig(dsn string) Config {
	conf := Config{}
	conf.parse(dsn)
	return conf
}

func (c Config) Dsn() string {
	return c.Data().String()
}

func (c Config) Data() *url.URL {
	return c.data
}
func (c Config) parse(dsn string) Config {
	data, err := url.Parse(dsn)
	if err != nil {
		data = nil
	}
	c.data = data
	return c
}
