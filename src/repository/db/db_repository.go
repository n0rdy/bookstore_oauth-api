package db

import (
	"bookstore_oauth-api/src/domain/access_token"
	"bookstore_oauth-api/src/domain/errors"
)

type DbRepository interface {
	GetById(string) (*access_token.AccessToken, *errors.RestErr)
}

type dbRepository struct{}

func NewRepository() DbRepository {
	return &dbRepository{}
}

func (r *dbRepository) GetById(string) (*access_token.AccessToken, *errors.RestErr) {
	return nil, errors.NewInternalServerError("db connection is not implemented yet")
}
