package users

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/avila-r/chat-hoster/internal/auth"
	"github.com/golang-jwt/jwt/v4"
)

type Service struct {
	r Repository
}

func NewService(r Repository) *Service {
	return &Service{r}
}

func (s *Service) Register(r *RegisterRequest) (*RegisterResponse, error) {
	hash, err := auth.EncryptPassword(r.Password)

	if err != nil {
		return nil, err
	}

	u := &User{
		Username: r.Username,
		Email:    r.Email,
		Password: hash,
	}

	if _, err := s.r.FindUserByEmail(u.Email); err == nil {
		return nil, errors.New("e-mail already registered")
	}

	// TODO: Error handling
	result, err := s.r.CreateUser(u)

	if err != nil {
		return nil, err
	}

	response := &RegisterResponse{
		ID:       int64(result.ID),
		Username: result.Username,
		Email:    result.Email,
	}

	return response, nil
}

func (s *Service) Login(r *LoginRequest) (*LoginResponse, error) {
	user, err := s.r.FindUserByEmail(r.Email)

	if err != nil {
		return nil, err
	}

	err = auth.CheckPassword(r.Password, user.Password)

	if err != nil {
		return nil, err
	}

	type Claims struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		jwt.RegisteredClaims
	}

	id := strconv.Itoa(int(user.ID))

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		ID:       id,
		Username: user.Username,

		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	var (
		secret = os.Getenv("a")
	)

	token, err := t.SignedString([]byte(secret))

	if err != nil {
		return nil, err
	}

	// TODO
	response := LoginResponse{
		Token:    token,
		ID:       id,
		Username: user.Username,
	}

	return &response, nil
}
