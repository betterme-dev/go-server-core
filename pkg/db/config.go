package db

import (
	"fmt"
	"strings"
)

type Config struct {
	Driver   string
	Username string
	Password string
	Host     string
	Port     int
	DbName   string
}

func NewConfig() *Config {
	return &Config{
		Driver: "mysql",
		Host:   "localhost",
		Port:   3306,
	}
}

// GetDsn builds database DSN string
func (dbc *Config) GetDsn() string {
	dsn := ""

	switch dbc.Driver {
	case "mysql":
		if dbc.Username != "" {
			dsn = dbc.Username
		}
		if dbc.Password != "" {
			dsn = fmt.Sprintf("%s:%s", dsn, dbc.Password)
		}
		dsn = fmt.Sprintf("%s@tcp(%s:%d)", dsn, dbc.Host, dbc.Port)
		if dbc.DbName != "" {
			dsn = fmt.Sprintf("%s/%s", dsn, dbc.DbName)
		}
	case "postgres":
		dsnParams := []string{
			"host=" + dbc.Host,
			"user=" + dbc.Username,
			"password=" + dbc.Password,
			"dbname=" + dbc.DbName,
			"sslmode=disable",
		}

		if dbc.Port != 0 {
			dsnParams = append(dsnParams, fmt.Sprintf("port=%d", dbc.Port))
		}

		dsn = strings.Join(dsnParams, " ")
	}

	return dsn
}
