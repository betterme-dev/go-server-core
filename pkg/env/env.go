package env

import (
	"database/sql"
	"github.com/betterme-dev/go-server-core/pkg/mq"
	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
)

var e *Env

type Env struct {
	Db       *sql.DB
	MqClient *mq.Client
	ES       *elasticsearch7.Client
}

func New() *Env {
	return &Env{}
}

func SetDB(db *sql.DB) {
	current().Db = db
}

func SetQueue(q *mq.Client) {
	current().MqClient = q
}

func SetElasticSearch(es *elasticsearch7.Client) {
	current().ES = es
}

func DB() *sql.DB {
	return current().Db
}

func Queue() *mq.Client {
	return current().MqClient
}

func ES() *elasticsearch7.Client {
	return current().ES
}

func current() *Env {
	if e == nil {
		e = New()
	}
	return e
}
