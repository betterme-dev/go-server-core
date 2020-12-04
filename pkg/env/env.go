package env

import (
	"database/sql"

	"github.com/betterme-dev/go-server-core/pkg/elasticsearch"
	"github.com/betterme-dev/go-server-core/pkg/mq"
	"github.com/betterme-dev/go-server-core/pkg/neo4j"

	"github.com/spf13/afero"
)

var e *Env

type Env struct {
	DB       *sql.DB
	DBSlave  *sql.DB
	MqClient *mq.Client
	Es       *elasticsearch.ES
	Fs       *afero.Fs
	Neo4j    *neo4j.Neo4j
}

func New() *Env {
	return &Env{}
}

func SetDB(db *sql.DB) {
	current().DB = db
}

func SetDBSlave(db *sql.DB) {
	current().DBSlave = db
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

func SetNeo4j(n *neo4j.Neo4j) {
	current().Neo4j = n
}

func DB() *sql.DB {
	return current().DB
}

func DBSlave() *sql.DB {
	if current().DBSlave != nil {
		return current().DBSlave
	}
	return current().DB
}

func Queue() *mq.Client {
	return current().MqClient
}

func ES() *elasticsearch.ES {
	return current().Es
}

func Neo4j() *neo4j.Neo4j {
	return current().Neo4j
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
