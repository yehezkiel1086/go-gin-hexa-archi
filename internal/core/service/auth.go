package service

import (
	"context"

	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/config"
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

func (as *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	// check email
	user, err := as.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	// check password
	if err := util.CheckPassword(user.Password, password); err != nil {
		return "", err
	}

	// generate jwt token
	return util.GenerateJWT(as.conf, user)
}
