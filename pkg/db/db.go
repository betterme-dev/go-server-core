package db

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/betterme-dev/go-server-core/pkg/env"
	_ "github.com/go-sql-driver/mysql"
	"github.com/huandu/go-sqlbuilder"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	preparedStatements      = map[string]*sql.Stmt{}
	preparedStatementsMutex = sync.RWMutex{}
)

func NewConnection() (*sql.DB, error) {
	dbConfig := NewConfig()
	dbConfig.Driver = viper.GetString("DB_DRIVER")
	dbConfig.Username = viper.GetString("DB_USERNAME")
	dbConfig.Password = viper.GetString("DB_PASSWORD")
	dbConfig.DbName = viper.GetString("DB_NAME")
	dbConfig.Host = viper.GetString("DB_HOST")

	port := viper.GetInt("DB_PORT")
	if port != 0 {
		dbConfig.Port = port
	}

	dbConn, err := NewWithConfig(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %s", err)
	}
	return dbConn, nil
}

func NewWithConfig(config *Config) (*sql.DB, error) {
	log.Infof("Connecting to %s database as %s to %s:%d/%s", config.Driver, config.Username, config.Host, config.Port, config.DbName)

	db, err := sql.Open(config.Driver, config.GetDsn())
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func ExecuteUpsert(qb *InsertBuilderWithDuplicateKeyUpdate) (id int64, err error) {
	res, err := execQuery(qb)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func ExecuteUpdate(qb *sqlbuilder.UpdateBuilder) (updated int64, err error) {
	res, err := execQuery(qb)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func FindOne(q *sqlbuilder.SelectBuilder, r interface{}) (found bool, err error) {
	query, args := q.Build()
	row := env.DB().QueryRow(query, args...)
	err = row.Scan(r)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// PrepareStatement creates a statement if it does not exists (and saves it to storage),
// otherwise takes it from storage.
func PrepareStatement(db *sql.DB, query string) (*sql.Stmt, error) {
	if stmt, exists := getPreparedStatementIfExists(query); exists {
		return stmt, nil
	}

	return prepareAndSave(db, query)
}

func getPreparedStatementIfExists(query string) (*sql.Stmt, bool) {
	preparedStatementsMutex.RLock()
	defer preparedStatementsMutex.RUnlock()

	stmt, exists := preparedStatements[query]
	return stmt, exists
}

func prepareAndSave(db *sql.DB, query string) (*sql.Stmt, error) {
	preparedStatementsMutex.Lock()
	defer preparedStatementsMutex.Unlock()

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}

	preparedStatements[query] = stmt
	return stmt, nil
}

func execQuery(qb sqlbuilder.Builder) (sql.Result, error) {
	query, args := qb.Build()
	log.Infof("Query: %s, args: %v", query, args)
	stmt, err := PrepareStatement(env.DB(), query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare DB statement")
	}
	return stmt.Exec(args...)
}
