package user

import (
	"github.com/betterme-dev/go-server-core/pkg/db"
	"github.com/doug-martin/goqu/v9"

	log "github.com/sirupsen/logrus"
)

const (
	TableNameUserSession = "user_session"
)

type SessionRepository struct {
	table string
	db    *goqu.Database
}

func NewSessionRepository() SessionRepository {
	return SessionRepository{
		table: TableNameUserSession,
		db:    db.Goqu(),
	}
}

func (sr SessionRepository) SessionByAuthKey(authKey string) (*Session, error) {
	log.Debugf("Search user by authKey '%s' in table %s", authKey, TableNameUserSession)
	var ses Session
	found, err := sr.db.
		From(sr.table).
		Select("id", "user_id", "expires_at").
		Where(
			goqu.C("auth_key").Eq(authKey),
			goqu.C("is_deleted").Eq(false),
		).
		Limit(1).
		ScanStruct(&ses)
	if err != nil {
		return nil, err
	}
	if !found {
		log.Debugf("user session not found in db with key %s", authKey)
		return nil, nil
	}
	return &ses, nil
}
