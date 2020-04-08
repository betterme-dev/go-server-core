package env

import (
	"database/sql"
	"github.com/betterme-dev/go-server-core/pkg/mq"
)

var e *Env

type Env struct {
	Db       *sql.DB
	MqClient *mq.Client
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

func DB() *sql.DB {
	return current().Db
}

func Queue() *mq.Client {
	return current().MqClient
}

func current() *Env {
	if e == nil {
		e = New()
	}
	return e
}
