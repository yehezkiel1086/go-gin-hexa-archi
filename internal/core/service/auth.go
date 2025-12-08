package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yehezkiel1086/go-gin-hexa-inventory/internal/adapter/config"
	"github.com/yehezkiel1086/go-gin-hexa-inventory/internal/core/domain"
	"github.com/yehezkiel1086/go-gin-hexa-inventory/internal/core/port"
	"github.com/yehezkiel1086/go-gin-hexa-inventory/internal/core/util"
)

type AuthService struct {
	conf *config.JWT
	repo port.UserRepository
}

func NewAuthService(conf *config.JWT, repo port.UserRepository) *AuthService {
	return &AuthService{
		conf: conf,
		repo: repo,
	}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	// check email
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	// check password
	if err := util.CheckPassword(user.Password, password); err != nil {
		return "", err
	}

	// create jwt
	mySigningKey := []byte(s.conf.Secret)

	duration, err := strconv.Atoi(s.conf.Duration)
	if err != nil {
		return "", err
	}

	// Create claims with multiple fields populated
	claims := domain.JWTClaims{
		Name: user.Name,
		Email: user.Email,
		Role: fmt.Sprint(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(duration) * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}

	return ss, nil
}
