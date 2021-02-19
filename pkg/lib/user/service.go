package user

type Service struct {
	Repository        Repository
	SessionRepository SessionRepository
}

func NewService() Service {
	return Service{
		Repository:        NewRepository(),
		SessionRepository: NewSessionRepository(),
	}
}

func (s Service) UserByID(id int) (*User, error) {
	return s.Repository.ByID(id)
}

func (s Service) UserByAuthToken(token string) (*User, error) {
	return s.Repository.ByAuthKey(token)
}

func (s Service) SessionByAuthToken(token string) (*Session, error) {
	return s.SessionRepository.SessionByAuthKey(token)
}
