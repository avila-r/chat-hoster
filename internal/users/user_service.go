package users

type Service struct {
	r Repository
}

func NewService(r Repository) *Service {
	return &Service{r}
}

func (s *Service) Register(r *RegisterRequest) (*RegisterResponse, error) {
	return nil, nil
}

func (s *Service) Login(r *LoginRequest) (*LoginResponse, error) {
	return nil, nil
}
