package service

import (
	"context"

	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/config"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/domain"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/port"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/util"
)

type AuthService struct {
	conf *config.JWT
	userRepo port.UserRepository
}

func NewAuthService(conf *config.JWT, userRepo port.UserRepository) *AuthService {
	return &AuthService{
		conf,
		userRepo,
	}
}

func (as *AuthService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := as.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}

	if err := util.CompareHashedPwd(user.Password, password); err != nil {
		return "", "", err
	}

	// generate jwt token
	refreshToken, err := util.GenerateJWTToken(as.conf, user, "refresh")
	if err != nil {
		return "", "", err
	}

	accessToken, err := util.GenerateJWTToken(as.conf, user, "access")
	if err != nil {
		return "", "", err
	}

	return refreshToken, accessToken, nil
}

func (as *AuthService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return as.userRepo.GetUserByEmail(ctx, email)
}
