package mq

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

const (
	defaultDriver          = "amqp"
	defaultHost            = "localhost"
	defaultPort            = "5672"
	paramConnectionTimeout = "connection_timeout"
)

type Config struct {
	Driver   string
	Username string
	Password string
	Host     string
	Port     string
	Vhost    string
	Params   map[string]string
}

func NewConfig() Config {
	return loadDefaultConfig()
}

func NewConfigFromDsn(dsn string) (Config, error) {
	conf, err := parseDsn(dsn)
	if err != nil {
		return Config{}, err
	}
	return *conf, nil
}

func (c Config) Dsn() string {
	params := ""
	if len(c.Params) != 0 {
		var pairs []string
		for k, v := range c.Params {
			pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
		}
		params = strings.Join(pairs, "&")
	}
	// avoid double-escape
	c.Username, _ = url.PathUnescape(c.Username)
	c.Password, _ = url.PathUnescape(c.Password)
	// compose
	dsn := fmt.Sprintf("%s://%s:%s@%s:%s/%s?%s",
		c.Driver,
		c.UsernameEscaped(),
		c.PasswordEscaped(),
		c.Host,
		c.Port,
		c.Vhost,
		params,
	)
	return dsn
}

func (c Config) UsernameEscaped() string {
	return url.PathEscape(c.Username)
}
func (c Config) PasswordEscaped() string {
	return url.PathEscape(c.Password)
}
func (c Config) VhostEscaped() string {
	return url.PathEscape(c.Vhost)
}
func (c Config) Param(name string) string {
	val, ok := c.Params[name]
	if !ok {
		return ""
	}
	return val
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
	// unescape auth data
	usr, err := url.PathUnescape(data.User.Username())
	if err != nil {
		usr = data.User.Username()
	}
	p, _ := data.User.Password()
	pass, err := url.PathUnescape(p)
	if err != nil {
		pass = p
	}
	// set data
	conf.Driver = data.Scheme
	conf.Username = usr
	conf.Password = pass
	conf.Host = data.Hostname()
	conf.Port = data.Port()
	conf.Vhost = strings.ReplaceAll(data.Path, "//", "/")
	params := make(map[string]string, len(data.Query()))
	for k, v := range data.Query() {
		params[k], _ = url.PathUnescape(v[0])
	}
	conf.Params = params
	return &conf, nil
}

func loadDefaultConfig() Config {
	var conf Config
	conf.Driver = defaultDriver
	conf.Host = defaultHost
	conf.Port = defaultPort
	conf.Params = make(map[string]string)
	return conf
}
