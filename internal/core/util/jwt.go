package util

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/adapter/config"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/domain"
)

func GenerateJWTToken(conf *config.JWT, user *domain.User, tokenType string) (string, error) {
	var mySigningKey []byte
	var duration int
	var err error
	var expiry *jwt.NumericDate

	if tokenType == "refresh" {
		mySigningKey = []byte(conf.RefreshTokenSecret)
		duration, err = strconv.Atoi(conf.RefreshTokenDuration)
		if err != nil {
			return "", err
		}
		expiry = jwt.NewNumericDate(time.Now().Add(time.Duration(duration) * time.Hour * 24))
	} else {
		mySigningKey = []byte(conf.AccessTokenSecret)
		duration, err = strconv.Atoi(conf.AccessTokenDuration)
		if err != nil {
			return "", err
		}
		expiry = jwt.NewNumericDate(time.Now().Add(time.Duration(duration) * time.Second))
	}

	// create claims
	claims := domain.JWTClaims{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expiry,
		},
	}

	// generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}

	return ss, nil
}

func ParseToken(tokenStr string, conf *config.JWT, tokenType string) (*domain.JWTClaims, error) {
	// parse token
	token, err := jwt.ParseWithClaims(tokenStr, &domain.JWTClaims{}, func(token *jwt.Token) (any, error) {
		if tokenType == "refresh" {
			return []byte(conf.RefreshTokenSecret), nil
		}

		return []byte(conf.AccessTokenSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	// extract claims
	claims, ok := token.Claims.(*domain.JWTClaims)
	if !ok {
		return nil, domain.ErrUnauthorized
	}

	return claims, nil
}
