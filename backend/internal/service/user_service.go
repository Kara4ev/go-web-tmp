package service

import (
	"context"

	"github.com/Kara4ev/go-web-tmp/internal/model"
	"github.com/Kara4ev/go-web-tmp/internal/model/apperrors"
	"github.com/Kara4ev/go-web-tmp/pkg/logger"
	"github.com/google/uuid"
)

type userService struct {
	UserRepository model.UserRepository
}

type USConfig struct {
	UserRepository model.UserRepository
}

func NewUserServices(c *USConfig) model.UserService {
	return &userService{
		UserRepository: c.UserRepository,
	}
}

func (s userService) Get(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	u, err := s.UserRepository.FindByID(ctx, uid)
	return u, err
}

func (s userService) Signup(ctx context.Context, u *model.User) error {

	pw, err := hashPassword(u.Password)

	if err != nil {
		logger.Warn("unable to signup user from email: %v", u.Email)
		return apperrors.NewInternal()
	}

	u.Password = pw

	if err := s.UserRepository.Create(ctx, u); err != nil {
		return err
	}

	return nil
}

func (s userService) Signin(ctx context.Context, u *model.User) error {
	uFetched, err := s.UserRepository.FindByEmail(ctx, u.Email)
	errAuthorization := apperrors.NewAuthorization("invalid email and password combination")
	if err != nil {
		logger.Warn("user search error by mail: %s, err: %v", u.Email, err)
		return errAuthorization
	}

	match, err := comparePassword(uFetched.Password, u.Password)

	if err != nil {
		logger.Error("error compare password, user email: %s", u.Email)
		return apperrors.NewInternal()
	}

	if !match {
		logger.Warn("invalid password, user email: %s", u.Email)
		return errAuthorization
	}

	*u = *uFetched
	return nil

}

func (s *userService) UpdateDetails(ctx context.Context, u *model.User) error {
	return s.UserRepository.Update(ctx, u)
}
