package auth

import (
	"errors"
	"time"

	"github.com/betterme-dev/go-server-core/pkg/lib/user"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	user user.User
}

func NewService() Service {
	return Service{
		user: user.User{},
	}
}

func (s Service) UserID() int {
	if s.user.ID != 0 {
		return s.user.ID
	}
	return 0
}

func (s Service) User() (*user.User, error) {
	if s.user.ID == 0 {
		return nil, errors.New("user not authorized")
	}
	return &s.user, nil
}

func (s *Service) AuthByBearerToken(token string) bool {
	usrSrv := user.NewService()
	now := int(time.Now().Unix())
	// check user data (deprecated)
	usr, err := usrSrv.UserByAuthToken(token)
	if err == nil && usr != nil && usr.AuthKeyExpires > now {
		return s.prepareUser(usr)
	}
	// check user-session data
	usrSes, err := usrSrv.SessionByAuthToken(token)
	if err != nil {
		log.Error(err)
		return false
	}
	if usrSes == nil {
		log.Debugf("Seesion not found")
		return false
	}
	if usrSes.ExpiresAt <= now {
		log.Debugf("Seesion token is expired")
		return false
	}
	usr, err = usrSrv.UserByID(usrSes.UserID)
	if err != nil {
		log.Error(err)
		return false
	}
	return s.prepareUser(usr)
}

func (s *Service) prepareUser(usr *user.User) bool {
	if usr == nil {
		return false
	}
	s.user = user.User{
		ID:             usr.ID,
		AuthKeyExpires: usr.AuthKeyExpires,
	}
	log.Debug("Logged in userID: ", usr.ID)

	return true
}
