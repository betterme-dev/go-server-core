package db

import (
	"github.com/betterme-dev/go-server-core/pkg/env"
	"github.com/doug-martin/goqu/v9"
	// import the dialect
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
)

var g *goqu.Database

func Goqu() *goqu.Database {
	if g == nil {
		g = goqu.Dialect("mysql").DB(env.DB())
	}
	return g
}
