package db

import "fmt"

type Config struct {
	Username string
	Password string
	Host     string
	Port     int
	DbName   string
}

func NewConfig() *Config {
	return &Config{
		Host: "localhost",
		Port: 3306,
	}
}

// GetDsn builds database DSN string
func (dbc *Config) GetDsn() string {
	dsn := ""
	if dbc.Username != "" {
		dsn = fmt.Sprintf("%s", dbc.Username)
	}
	if dbc.Password != "" {
		dsn = fmt.Sprintf("%s:%s", dsn, dbc.Password)
	}
	dsn = fmt.Sprintf("%s@tcp(%s:%d)", dsn, dbc.Host, dbc.Port)
	if dbc.DbName != "" {
		dsn = fmt.Sprintf("%s/%s", dsn, dbc.DbName)
	}

	return dsn
}
