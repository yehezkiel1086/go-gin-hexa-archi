package service

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/domain"
	"github.com/yehezkiel1086/go-gin-hexa-archi/internal/core/util"
	"github.com/yehezkiel1086/go-gin-hexa-archi/mocks"
)

type registerTestedInput struct {
	user *domain.User
}

type registerExpectedOutput struct {
	user *domain.UserResponse
	err  error
}

func TestUserService_Register(t *testing.T) {
	ctx := context.Background()

	// arange
	userName := gofakeit.Name()
	userEmail := gofakeit.Email()
	userPassword := gofakeit.Password(true, true, true, true, false, 8)

	userInput := &domain.User{
		Name:     userName,
		Email:    userEmail,
		Password: userPassword,
	}

	userResponse := &domain.UserResponse{
		Name:  userName,
		Email: userEmail,
	}

	cacheKey := util.GenerateCacheKey("user", 0)
	ttl := time.Duration(0)

	testCases := []struct {
		desc  string
		mocks func(
			userRepo *mocks.UserRepository,
			cache *mocks.CacheRepository,
		)
		input    registerTestedInput
		expected registerExpectedOutput
	}{
		{
			desc: "Success",
			mocks: func(
				userRepo *mocks.UserRepository,
				cache *mocks.CacheRepository,
			) {
				userRepo.
					On("CreateUser", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
						return u.Email == userInput.Email && u.Name == userInput.Name
					})).
					Return(userResponse, nil).
					Once()

				cache.
					On("Set", mock.Anything, cacheKey, mock.Anything, ttl).
					Return(nil).
					Once()

				cache.
					On("DeleteByPrefix", mock.Anything, "users:*").
					Return(nil).
					Once()
			},
			input: registerTestedInput{
				user: userInput,
			},
			expected: registerExpectedOutput{
				user: userResponse,
				err:  nil,
			},
		},
		{
			desc: "Fail_InternalError",
			mocks: func(
				userRepo *mocks.UserRepository,
				cache *mocks.CacheRepository,
			) {
				cache.
					On("Set", mock.Anything, cacheKey, mock.Anything, ttl).
					Return(nil).
					Once()

				cache.
					On("DeleteByPrefix", mock.Anything, "users:*").
					Return(nil).
					Once()

				userRepo.
					On("CreateUser", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
						return u.Email == userInput.Email && u.Name == userInput.Name
					})).
					Return(nil, domain.ErrInternal).
					Once()
			},
			input: registerTestedInput{
				user: userInput,
			},
			expected: registerExpectedOutput{
				user: nil,
				err:  domain.ErrInternal,
			},
		},
		{
			desc: "Fail_DuplicateData",
			mocks: func(
				userRepo *mocks.UserRepository,
				cache *mocks.CacheRepository,
			) {
				cache.
					On("Set", mock.Anything, cacheKey, mock.Anything, ttl).
					Return(nil).
					Once()

				cache.
					On("DeleteByPrefix", mock.Anything, "users:*").
					Return(nil).
					Once()

				userRepo.
					On("CreateUser", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
						return u.Email == userInput.Email && u.Name == userInput.Name
					})).
					Return(nil, domain.ErrConflictingData).
					Once()
			},
			input: registerTestedInput{
				user: userInput,
			},
			expected: registerExpectedOutput{
				user: nil,
				err:  domain.ErrConflictingData,
			},
		},
		{
			desc: "Fail_SetCache",
			mocks: func(
				userRepo *mocks.UserRepository,
				cache *mocks.CacheRepository,
			) {
				cache.
					On("Set", mock.Anything, cacheKey, mock.Anything, ttl).
					Return(domain.ErrInternal).
					Once()
			},
			input: registerTestedInput{
				user: userInput,
			},
			expected: registerExpectedOutput{
				user: nil,
				err:  domain.ErrInternal,
			},
		},
		{
			desc: "Fail_DeleteCacheByPrefix",
			mocks: func(
				userRepo *mocks.UserRepository,
				cache *mocks.CacheRepository,
			) {
				cache.
					On("Set", mock.Anything, cacheKey, mock.Anything, ttl).
					Return(nil).
					Once()

				cache.
					On("DeleteByPrefix", mock.Anything, "users:*").
					Return(domain.ErrInternal).
					Once()
			},
			input: registerTestedInput{
				user: userInput,
			},
			expected: registerExpectedOutput{
				user: nil,
				err:  domain.ErrInternal,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			userRepo := new(mocks.UserRepository)
			cache := new(mocks.CacheRepository)

			tc.mocks(userRepo, cache)

			userService := NewUserService(userRepo, cache)

			// Clone input to avoid side effects (hashing) on the shared struct
			input := &domain.User{
				Name:     tc.input.user.Name,
				Email:    tc.input.user.Email,
				Password: tc.input.user.Password,
			}
			user, err := userService.RegisterUser(ctx, input)

			assert.Equal(t, tc.expected.err, err)
			assert.Equal(t, tc.expected.user, user)

			userRepo.AssertExpectations(t)
			cache.AssertExpectations(t)
		})
	}
}
