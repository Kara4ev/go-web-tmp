package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/Kara4ev/go-web-tmp/internal/model"
	"github.com/Kara4ev/go-web-tmp/internal/model/apperrors"
	"github.com/Kara4ev/go-web-tmp/internal/model/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGet(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUserResp := &model.User{
			UID:   uid,
			Email: "bob@bob.com",
			Name:  "Bobby Bobson",
		}

		mockUserRepository := new(mocks.MockUserRepository)

		us := NewUserServices(&USConfig{
			UserRepository: mockUserRepository,
		})

		mockUserRepository.On("FindByID", mock.Anything, uid).Return(mockUserResp, nil)

		ctx := context.TODO()
		u, err := us.Get(ctx, uid)

		assert.NoError(t, err)
		assert.Equal(t, u, mockUserResp)
		mockUserRepository.AssertExpectations(t)

	})

	t.Run("Error", func(t *testing.T) {

		mockUserRepository := new(mocks.MockUserRepository)

		us := NewUserServices(&USConfig{
			UserRepository: mockUserRepository,
		})

		mockUserRepository.On("FindByID", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("some error down call chain"))

		ctx := context.TODO()
		uid, _ := uuid.NewRandom()
		u, err := us.Get(ctx, uid)

		assert.Error(t, err)
		assert.Nil(t, u)

		mockUserRepository.AssertExpectations(t)

	})
}

func TestSignup(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		uid, _ := uuid.NewRandom()

		mockUser := &model.User{
			Email:    "correct@email.com",
			Password: "correct-password",
		}

		mockUserRepository := new(mocks.MockUserRepository)
		us := NewUserServices(&USConfig{
			UserRepository: mockUserRepository,
		})

		mockUserRepository.
			On("Create", mock.Anything, mockUser).
			Run(func(args mock.Arguments) {
				userArg := args.Get(1).(*model.User)
				userArg.UID = uid
			}).
			Return(nil)

		ctx := context.TODO()
		err := us.Signup(ctx, mockUser)

		assert.NoError(t, err)
		assert.Equal(t, uid, mockUser.UID)
		mockUserRepository.AssertExpectations(t)

	})

	t.Run("Error", func(t *testing.T) {
		mockUser := &model.User{
			Email:    "correct@email.com",
			Password: "correct-password",
		}

		mockError := apperrors.NewConflict("email", mockUser.Email)

		mockUserRepository := new(mocks.MockUserRepository)
		us := NewUserServices(&USConfig{
			UserRepository: mockUserRepository,
		})

		mockUserRepository.
			On("Create", mock.Anything, mockUser).
			Return(mockError)

		ctx := context.TODO()
		err := us.Signup(ctx, mockUser)

		assert.Error(t, err)
		assert.Equal(t, err, mockError)
		assert.EqualError(t, err, mockError.Error())
		mockUserRepository.AssertExpectations(t)
	})
}

func TestSignin(t *testing.T) {

	mockEmail := "correct@email.com"
	mockCorrectHashPassword := "2232269800b344a31f9a5b5ca6c91775dc30c5d856d1a89c011076c6437236a5.52fdfc072182654f163f5f0f9a621d729566c74d10037c4d7bbb0407d1e2c649"
	mockInvalidHashPassword := "incorrect.hash"
	mockCorrectPassword := "correct-password"
	mockIncorrectPassword := "incorrect-password"

	t.Run("Success", func(t *testing.T) {

		mockUser := &model.User{
			Email:    mockEmail,
			Password: mockCorrectPassword,
		}

		mockFetchUser := &model.User{
			Email:    mockEmail,
			Password: mockCorrectHashPassword,
		}

		mockUserRepository := new(mocks.MockUserRepository)

		mockUserRepository.
			On("FindByEmail", mock.AnythingOfType("*context.emptyCtx"), mockEmail).
			Return(mockFetchUser, nil)

		us := NewUserServices(&USConfig{
			UserRepository: mockUserRepository,
		})

		ctx := context.TODO()
		err := us.Signin(ctx, mockUser)
		assert.NoError(t, err)
		mockUserRepository.AssertExpectations(t)
	})

	t.Run("Error -> no user email repository", func(t *testing.T) {

		errAuthorization := apperrors.NewAuthorization("some error")

		mockUser := &model.User{
			Email:    mockEmail,
			Password: mockCorrectPassword,
		}

		mockUserRepository := new(mocks.MockUserRepository)

		mockUserRepository.
			On("FindByEmail", mock.AnythingOfType("*context.emptyCtx"), mockEmail).
			Return(nil, fmt.Errorf("some error"))

		us := NewUserServices(&USConfig{
			UserRepository: mockUserRepository,
		})

		ctx := context.TODO()
		err := us.Signin(ctx, mockUser)
		assert.IsType(t, errAuthorization, err)
		mockUserRepository.AssertExpectations(t)
	})

	t.Run("Error -> incorrect user password", func(t *testing.T) {

		errAuthorization := apperrors.NewAuthorization("some error")

		mockUser := &model.User{
			Email:    mockEmail,
			Password: mockIncorrectPassword,
		}

		mockFetchUser := &model.User{
			Email:    mockEmail,
			Password: mockCorrectHashPassword,
		}

		mockUserRepository := new(mocks.MockUserRepository)

		mockUserRepository.
			On("FindByEmail", mock.AnythingOfType("*context.emptyCtx"), mockEmail).
			Return(mockFetchUser, nil)

		us := NewUserServices(&USConfig{
			UserRepository: mockUserRepository,
		})

		ctx := context.TODO()
		err := us.Signin(ctx, mockUser)
		assert.IsType(t, errAuthorization, err)
		mockUserRepository.AssertExpectations(t)
	})

	t.Run("Error -> incorrect hash password", func(t *testing.T) {

		errInternal := apperrors.NewInternal()
		mockUser := &model.User{
			Email:    mockEmail,
			Password: mockCorrectPassword,
		}

		mockFetchUser := &model.User{
			Email:    mockEmail,
			Password: mockInvalidHashPassword,
		}

		mockUserRepository := new(mocks.MockUserRepository)

		mockUserRepository.
			On("FindByEmail", mock.AnythingOfType("*context.emptyCtx"), mockEmail).
			Return(mockFetchUser, nil)

		us := NewUserServices(&USConfig{
			UserRepository: mockUserRepository,
		})

		ctx := context.TODO()
		err := us.Signin(ctx, mockUser)
		assert.IsType(t, errInternal, err)
		mockUserRepository.AssertExpectations(t)
	})

}
