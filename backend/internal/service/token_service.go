package service

import (
	"context"
	"crypto/rsa"

	"github.com/Kara4ev/go-web-tmp/internal/model"
	"github.com/Kara4ev/go-web-tmp/internal/model/apperrors"
	"github.com/Kara4ev/go-web-tmp/pkg/logger"
	"github.com/google/uuid"
)

type tokenService struct {
	TokenRepository       model.TokenRepository
	PrivKey               *rsa.PrivateKey
	PubKey                *rsa.PublicKey
	RefreshSecret         string
	IDExpirationSecs      int64
	RefrashExpirationSecs int64
}

type TSConfig struct {
	TokenRepository       model.TokenRepository
	PrivKey               *rsa.PrivateKey
	PubKey                *rsa.PublicKey
	RefreshSecret         string
	IDExpirationSecs      int64
	RefrashExpirationSecs int64
}

func NewTokenService(c *TSConfig) model.TokenService {
	return &tokenService{
		TokenRepository:       c.TokenRepository,
		PrivKey:               c.PrivKey,
		PubKey:                c.PubKey,
		RefreshSecret:         c.RefreshSecret,
		IDExpirationSecs:      c.IDExpirationSecs,
		RefrashExpirationSecs: c.RefrashExpirationSecs,
	}
}

func (s *tokenService) NewPairFromUser(ctx context.Context, u *model.User, prevTokenID string) (*model.TokenPair, error) {

	if prevTokenID != "" {
		if err := s.TokenRepository.DeleteRefreshToken(ctx, u.UID.String(), prevTokenID); err != nil {
			logger.Warn("error delete repository prev token: %v, for uid: %v, error: %v", prevTokenID, u.UID, err.Error())
			return nil, err
		}
	}

	idToken, err := generateIDToken(u, s.PrivKey, s.IDExpirationSecs)

	if err != nil {
		logger.Warn("error generating id token for uid: %v, error: %v", u.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	refreshToken, err := generateRefrashToken(u.UID, s.RefreshSecret, s.RefrashExpirationSecs)

	if err != nil {
		logger.Warn("error generating refresh token for uid: %v, error: %v", u.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	if err := s.TokenRepository.SetRefreshToken(ctx, u.UID.String(), refreshToken.ID.String(), refreshToken.ExpiresIn); err != nil {
		logger.Warn("error set repository refresh token: %v,  for uid: %v, error: %v", refreshToken.ID, u.UID, err.Error())
		return nil, apperrors.NewInternal()
	}

	return &model.TokenPair{
		IDToken:      model.IDToken{SS: idToken},
		RefreshToken: model.RefreshToken{ID: refreshToken.ID, UID: u.UID, SS: refreshToken.SS},
	}, nil
}

func (s *tokenService) ValidateIDToken(tokenString string) (*model.User, error) {

	claims, err := validateIDToken(tokenString, s.PubKey)
	if err != nil {
		logger.Warn("id token is invalid: %s, err: %v", tokenString, err)
		return nil, apperrors.NewAuthorization("unable to veryfy user from id token")
	}

	return claims.User, nil
}

func (s *tokenService) ValidateRefreshToken(tokenString string) (*model.RefreshToken, error) {
	claims, err := validateRefreshToken(tokenString, s.RefreshSecret)

	if err != nil {
		logger.Warn("refresh token is invalid: %s, err: %v", tokenString, err)
		return nil, apperrors.NewAuthorization("unable to veryfy user from id token")
	}

	tokensUUID, err := uuid.Parse(claims.Id)
	if err != nil {
		logger.Warn("claims ID could not be parsed as uuid: %s, err: %v", claims.Id, err)
		return nil, apperrors.NewAuthorization("unable to veryfy user from id token")
	}

	return &model.RefreshToken{
		ID:  tokensUUID,
		SS:  tokenString,
		UID: claims.UID,
	}, nil
}

func (s *tokenService) Signout(ctx context.Context, uid uuid.UUID) error {
	return s.TokenRepository.DeleteUserRefreshToken(ctx, uid.String())
}
