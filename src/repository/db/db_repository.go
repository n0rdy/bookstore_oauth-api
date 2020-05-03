package db

import (
	"github.com/gocql/gocql"
	"github.com/n0rdy/bookstore_oauth-api/src/clients/cassandra"
	"github.com/n0rdy/bookstore_oauth-api/src/domain/access_token"
	"github.com/n0rdy/bookstore_oauth-api/src/domain/errors"
)

const (
	queryGetAccessToken       = "SELECT access_token, user_id, client_id, expires FROM access_tokens WHERE access_token = ?;"
	queryCreateAccessToken    = "INSERT INTO access_tokens (access_token, user_id, client_id, expires) VALUES (?, ?, ?, ?);"
	queryUpdateExpirationTime = "UPDATE access_tokens SET expires = ? WHERE access_token = ?;"
)

type DbRepository interface {
	GetById(string) (*access_token.AccessToken, *errors.RestErr)
	Create(access_token.AccessToken) *errors.RestErr
	UpdateExpirationTime(access_token.AccessToken) *errors.RestErr
}

type dbRepository struct{}

func NewRepository() DbRepository {
	return &dbRepository{}
}

func (r *dbRepository) GetById(id string) (*access_token.AccessToken, *errors.RestErr) {
	session := cassandra.GetSession()
	var accessToken access_token.AccessToken

	if err := session.Query(queryGetAccessToken, id).Scan(
		&accessToken.AccessToken,
		&accessToken.UserId,
		&accessToken.ClientId,
		&accessToken.Expires,
	); err != nil {
		if err == gocql.ErrNotFound {
			return nil, errors.NewNotFoundError("no access token found with given id")
		}
		return nil, errors.NewInternalServerError(err.Error())
	}

	return &accessToken, nil
}

func (r *dbRepository) Create(token access_token.AccessToken) *errors.RestErr {
	session := cassandra.GetSession()

	if err := session.Query(queryCreateAccessToken,
		token.AccessToken,
		token.UserId,
		token.ClientId,
		token.Expires,
	).Exec(); err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	return nil
}

func (r *dbRepository) UpdateExpirationTime(token access_token.AccessToken) *errors.RestErr {
	session := cassandra.GetSession()

	if err := session.Query(queryUpdateExpirationTime,
		token.Expires,
		token.AccessToken,
	).Exec(); err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	return nil
}
