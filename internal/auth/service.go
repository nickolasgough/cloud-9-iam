package auth

type Service interface {
	CreateSession()
	ValidateSession(session string)
}

type service struct{}

func New() Service {
	return &service{}
}

func (s *service) CreateSession() {
}

func (s *service) ValidateSession(session string) {
}
