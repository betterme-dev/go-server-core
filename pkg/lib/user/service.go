package user

type Service struct {
	Repository Repository
}

func NewService() Service {
	return Service{Repository: NewRepository()}
}

func (s Service) GetUserByAuthToken(token string) (*User, error) {
	return s.Repository.GetByAuthKey(token)
}
