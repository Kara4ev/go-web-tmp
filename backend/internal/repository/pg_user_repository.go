package repository

import (
	"context"

	"github.com/Kara4ev/go-web-tmp/internal/model"
	"github.com/Kara4ev/go-web-tmp/internal/model/apperrors"
	"github.com/Kara4ev/go-web-tmp/pkg/logger"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type pgUserRepository struct {
	DB *sqlx.DB
}

func NewUserReposytory(db *sqlx.DB) model.UserRepository {
	return &pgUserRepository{
		DB: db,
	}
}

func (r *pgUserRepository) FindByID(ctx context.Context, uid uuid.UUID) (*model.User, error) {

	user := new(model.User)
	query := "SELECT * FROM users WHERE uid = $1"

	if err := r.DB.GetContext(ctx, user, query, uid); err != nil {
		logger.Warn("unable to get user with uid: %v, err: %v", uid.String(), err.Error())
		return nil, apperrors.NewNotFound("uid", uid.String())
	}
	return user, nil

}

func (r *pgUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {

	user := new(model.User)

	query := "SELECT * FROM users WHERE email=$1"

	if err := r.DB.GetContext(ctx, user, query, email); err != nil {
		logger.Warn("unable to get user with email addres: %v. Err: %v", email, err.Error())
		return nil, err
	}
	return user, nil

}

func (r *pgUserRepository) Create(ctx context.Context, u *model.User) error {
	query := "INSERT INTO users (email, password) VALUES ($1, $2) RETURNING *"
	if err := r.DB.GetContext(ctx, u, query, u.Email, u.Password); err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
			logger.Warn("cloud not create a user with email: %v , reason: %v", u.Email, err.Code.Name())
			return apperrors.NewConflict("email", u.Email)
		}

		logger.Warn("cloud not create a user with email: %v , reason: %v", u.Email, err)
		return apperrors.NewInternal()
	}
	return nil
}

func (r *pgUserRepository) Update(ctx context.Context, u *model.User) error {
	query := `
		UPDATE
			users
		SET 
			name=:name,
			email=:email
		WHERE
			uid=:uid
		RETURNING *;`

	nstmt, err := r.DB.PrepareNamedContext(ctx, query)
	if err != nil {
		logger.Warn("unable to prepare user update query: %v", err)
		return apperrors.NewInternal()
	}

	if err := nstmt.GetContext(ctx, u, u); err != nil {
		logger.Warn("Unable to prepare user update query: %v", err)
		return apperrors.NewInternal()
	}

	return nil

}
