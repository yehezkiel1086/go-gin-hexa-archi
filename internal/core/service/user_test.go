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

func TestUserService_GetUsers(t *testing.T) {
	ctx := context.Background()
	start := uint64(0)
	end := uint64(10)

	userResponse := domain.UserResponse{
		ID:    uint(gofakeit.Number(1, 100)),
		Name:  gofakeit.Name(),
		Email: gofakeit.Email(),
	}
	users := []domain.UserResponse{userResponse}

	params := util.GenerateCacheKeyParams(start, end)
	cacheKey := util.GenerateCacheKey("users", params)
	serializedUsers, _ := util.Serialize(users)

	testCases := []struct {
		desc     string
		mocks    func(*mocks.UserRepository, *mocks.CacheRepository)
		expected []domain.UserResponse
		err      error
	}{
		{
			desc: "Success_CacheHit",
			mocks: func(ur *mocks.UserRepository, cr *mocks.CacheRepository) {
				cr.On("Get", ctx, cacheKey).Return(serializedUsers, nil).Once()
			},
			expected: users,
			err:      nil,
		},
		{
			desc: "Success_CacheMiss",
			mocks: func(ur *mocks.UserRepository, cr *mocks.CacheRepository) {
				cr.On("Get", ctx, cacheKey).Return(nil, domain.ErrInternal).Once()
				ur.On("GetUsers", ctx, start, end).Return(users, nil).Once()
				cr.On("Set", ctx, cacheKey, mock.Anything, time.Duration(0)).Return(nil).Once()
			},
			expected: users,
			err:      nil,
		},
		{
			desc: "Fail_RepoError",
			mocks: func(ur *mocks.UserRepository, cr *mocks.CacheRepository) {
				cr.On("Get", ctx, cacheKey).Return(nil, domain.ErrInternal).Once()
				ur.On("GetUsers", ctx, start, end).Return(nil, domain.ErrInternal).Once()
			},
			expected: nil,
			err:      domain.ErrInternal,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ur := new(mocks.UserRepository)
			cr := new(mocks.CacheRepository)
			tc.mocks(ur, cr)

			s := NewUserService(ur, cr)
			res, err := s.GetUsers(ctx, start, end)

			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.expected, res)
			ur.AssertExpectations(t)
			cr.AssertExpectations(t)
		})
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	ctx := context.Background()
	id := uint(gofakeit.Number(1, 100))

	user := &domain.User{
		ID:    id,
		Name:  gofakeit.Name(),
		Email: gofakeit.Email(),
	}

	cacheKey := util.GenerateCacheKey("user", id)
	serializedUser, _ := util.Serialize(user)

	testCases := []struct {
		desc     string
		mocks    func(*mocks.UserRepository, *mocks.CacheRepository)
		expected *domain.User
		err      error
	}{
		{
			desc: "Success_CacheHit",
			mocks: func(ur *mocks.UserRepository, cr *mocks.CacheRepository) {
				cr.On("Get", ctx, cacheKey).Return(serializedUser, nil).Once()
			},
			expected: user,
			err:      nil,
		},
		{
			desc: "Success_CacheMiss",
			mocks: func(ur *mocks.UserRepository, cr *mocks.CacheRepository) {
				cr.On("Get", ctx, cacheKey).Return(nil, domain.ErrInternal).Once()
				ur.On("GetUserByID", ctx, id).Return(user, nil).Once()
				cr.On("Set", ctx, cacheKey, mock.Anything, time.Duration(0)).Return(nil).Once()
			},
			expected: user,
			err:      nil,
		},
		{
			desc: "Fail_RepoError",
			mocks: func(ur *mocks.UserRepository, cr *mocks.CacheRepository) {
				cr.On("Get", ctx, cacheKey).Return(nil, domain.ErrInternal).Once()
				ur.On("GetUserByID", ctx, id).Return(nil, domain.ErrInternal).Once()
			},
			expected: nil,
			err:      domain.ErrInternal,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ur := new(mocks.UserRepository)
			cr := new(mocks.CacheRepository)
			tc.mocks(ur, cr)

			s := NewUserService(ur, cr)
			res, err := s.GetUserByID(ctx, id)

			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.expected, res)
			ur.AssertExpectations(t)
			cr.AssertExpectations(t)
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	ctx := context.Background()
	id := uint(gofakeit.Number(1, 100))

	updateInput := &domain.User{
		Name: "New Name",
	}

	cacheKey := util.GenerateCacheKey("user", id)

	testCases := []struct {
		desc  string
		mocks func(*mocks.UserRepository, *mocks.CacheRepository, *domain.User)
		err   error
	}{
		{
			desc: "Success",
			mocks: func(ur *mocks.UserRepository, cr *mocks.CacheRepository, existing *domain.User) {
				serialized, _ := util.Serialize(existing)
				cr.On("Get", ctx, cacheKey).Return(serialized, nil).Once()
				// Note: Implementation calls GetUserByID even if cache hit
				ur.On("GetUserByID", ctx, id).Return(existing, nil).Once()
				ur.On("UpdateUser", ctx, mock.MatchedBy(func(u *domain.User) bool {
					return u.Name == updateInput.Name
				})).Return(existing, nil).Once()
				cr.On("Delete", ctx, cacheKey).Return(nil).Once()
				cr.On("Delete", ctx, "users:*").Return(nil).Once()
				cr.On("Set", ctx, cacheKey, mock.Anything, time.Duration(0)).Return(nil).Once()
			},
			err: nil,
		},
		{
			desc: "Fail_RepoUpdateError",
			mocks: func(ur *mocks.UserRepository, cr *mocks.CacheRepository, existing *domain.User) {
				cr.On("Get", ctx, cacheKey).Return(nil, domain.ErrInternal).Once()
				ur.On("GetUserByID", ctx, id).Return(existing, nil).Once()
				ur.On("UpdateUser", ctx, mock.Anything).Return(nil, domain.ErrInternal).Once()
			},
			err: domain.ErrInternal,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ur := new(mocks.UserRepository)
			cr := new(mocks.CacheRepository)

			// Create fresh existing user for each run to avoid side effects
			existingUser := &domain.User{
				ID:       id,
				Name:     gofakeit.Name(),
				Email:    gofakeit.Email(),
				Password: "oldpassword",
			}

			tc.mocks(ur, cr, existingUser)

			s := NewUserService(ur, cr)
			res, err := s.UpdateUser(ctx, id, updateInput)

			assert.Equal(t, tc.err, err)
			if tc.err == nil {
				assert.Equal(t, updateInput.Name, res.Name)
			}
			ur.AssertExpectations(t)
			cr.AssertExpectations(t)
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	ctx := context.Background()
	id := uint(gofakeit.Number(1, 100))
	cacheKey := util.GenerateCacheKey("user", id)
	deletedUser := &domain.User{ID: id}

	testCases := []struct {
		desc     string
		mocks    func(*mocks.UserRepository, *mocks.CacheRepository)
		expected *domain.User
		err      error
	}{
		{
			desc: "Success",
			mocks: func(ur *mocks.UserRepository, cr *mocks.CacheRepository) {
				cr.On("Delete", ctx, cacheKey).Return(nil).Once()
				cr.On("Delete", ctx, "users:*").Return(nil).Once()
				ur.On("DeleteUser", ctx, id).Return(deletedUser, nil).Once()
			},
			expected: deletedUser,
			err:      nil,
		},
		{
			desc: "Fail_CacheDeleteError",
			mocks: func(ur *mocks.UserRepository, cr *mocks.CacheRepository) {
				cr.On("Delete", ctx, cacheKey).Return(domain.ErrInternal).Once()
			},
			expected: nil,
			err:      domain.ErrInternal,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ur := new(mocks.UserRepository)
			cr := new(mocks.CacheRepository)
			tc.mocks(ur, cr)

			s := NewUserService(ur, cr)
			res, err := s.DeleteUser(ctx, id)

			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.expected, res)
			ur.AssertExpectations(t)
			cr.AssertExpectations(t)
		})
	}
}
