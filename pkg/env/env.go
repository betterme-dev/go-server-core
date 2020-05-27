package env

import (
	"database/sql"
	"github.com/betterme-dev/go-server-core/pkg/elasticsearch"
	"github.com/betterme-dev/go-server-core/pkg/mq"
	"github.com/spf13/afero"
)

var e *Env

type Env struct {
	Db       *sql.DB
	MqClient *mq.Client
	Es       *elasticsearch.ES
	Fs       *afero.Fs
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

func SetElasticSearch(es *elasticsearch.ES) {
	current().Es = es
}

func SetFS(fs *afero.Fs) {
	current().Fs = fs
}

func DB() *sql.DB {
	return current().Db
}

func Queue() *mq.Client {
	return current().MqClient
}

func ES() *elasticsearch.ES {
	return current().Es
}

func FS() *afero.Fs {
	return current().Fs
}

func current() *Env {
	if e == nil {
		e = New()
	}
	return e
}
