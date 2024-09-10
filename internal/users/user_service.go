package users

type Service struct {
	r Repository
}

func NewService(r Repository) *Service {
	return &Service{r}
}

func (s *Service) Register(r *RegisterRequest) (*RegisterResponse, error) {
	u := &User{
		Username: r.Username,
		Email:    r.Email,
		Password: r.Password,
	}


	// TODO: Error handling
	result, _ := s.r.CreateUser(u)

	response := &RegisterResponse{
		ID:       int64(result.ID),
		Username: result.Email,
		Email:    result.Password,
	}

	return response, nil
}

func (s *Service) Login(r *LoginRequest) (*LoginResponse, error) {
	// TODO
	response := LoginResponse{
		Token:    "access_token",
		ID:       1,
		Username: "username",
	}

	return &response, nil
}
