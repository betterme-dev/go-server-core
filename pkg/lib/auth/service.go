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
	usr, err := user.NewService().GetUserByAuthToken(token)
	if err != nil {
		log.Error(err)
		return false
	}
	if usr == nil {
		log.Debugf("User with token %s not found", token)
		return false
	}
	if usr.AuthKeyExpires <= time.Now().Unix() {
		log.Debugf("Token %s is expired", token)
		return false
	}
	s.user = user.User{
		ID:             usr.ID,
		AuthKeyExpires: usr.AuthKeyExpires,
	}
	log.Debugf("Logged in userID: %d", usr.ID)
	return true
}
