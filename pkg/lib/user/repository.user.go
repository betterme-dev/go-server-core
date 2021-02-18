package user

import (
	"github.com/betterme-dev/go-server-core/pkg/db"
	"github.com/doug-martin/goqu/v9"

	log "github.com/sirupsen/logrus"
)

const (
	TableNameUser = "user"
)

type Repository struct {
	table string
	db    *goqu.Database
}

func NewRepository() Repository {
	return Repository{
		table: TableNameUser,
		db:    db.Goqu(),
	}
}

func (r Repository) ByID(id int) (*User, error) {
	log.Debugf("Search user by user ID '%d' in table %s", id, r.table)
	var user User
	found, err := r.db.
		From(r.table).
		Select("id", "auth_key_expires").
		Where(
			goqu.C("id").Eq(id),
		).
		Limit(1).
		ScanStruct(&user)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}

	return &user, nil
}

func (r Repository) ByAuthKey(authKey string) (*User, error) {
	log.Debugf("Search user by authKey '%s' in table %s", authKey, r.table)
	var user User
	found, err := r.db.
		From(r.table).
		Select("id", "auth_key_expires").
		Where(
			goqu.C("auth_key").Eq(authKey),
		).
		Limit(1).
		ScanStruct(&user)
	if err != nil {
		return nil, err
	}
	if !found {
		log.Debugf("user not found in db with key %s", authKey)
		return nil, nil
	}
	return &user, nil
}
