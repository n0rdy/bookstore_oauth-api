package rest

import (
	"bookstore_oauth-api/src/domain/errors"
	"bookstore_oauth-api/src/domain/users"
	"encoding/json"
	"github.com/mercadolibre/golang-restclient/rest"
	"time"
)

var (
	usersRestClient = rest.RequestBuilder{
		BaseURL: "https://api.bookstore.com", // TODO change to a real one
		Timeout: 100 * time.Millisecond,
	}
)

type RestUsersRepository interface {
	LoginUser(string, string) (*users.User, *errors.RestErr)
}

type usersRepository struct{}

func NewRestUsersRepository() RestUsersRepository {
	return &usersRepository{}
}

func (ur *usersRepository) LoginUser(email string, password string) (*users.User, *errors.RestErr) {
	request := users.UserLoginRequest{
		Email:    email,
		Password: password,
	}
	response := usersRestClient.Post("/users/login", request)
	if response == nil || response.Response == nil {
		return nil, errors.NewInternalServerError("invalid rest client response on trying to login a user")
	}
	if response.StatusCode > 299 {
		var restErr errors.RestErr
		err := json.Unmarshal(response.Bytes(), &restErr)
		if err != nil {
			return nil, errors.NewInternalServerError("invalid error interface on trying to login a user")
		}
		return nil, &restErr
	}

	var user users.User
	err := json.Unmarshal(response.Bytes(), user)
	if err != nil {
		return nil, errors.NewInternalServerError("error on trying to unmarshal users response")
	}

	return &user, nil
}
