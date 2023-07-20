package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/iliyaisd/littlejohn/ljlib"
)

const mockPassword = ""

type APIKeyAuthorizer struct {
	userRepository UserRepository
}

type UserRepository interface {
	GetUserByUsername(username string) (*ljlib.User, error)
}

func NewAPIKeyAuthorizer(repository UserRepository) APIKeyAuthorizer {
	return APIKeyAuthorizer{
		userRepository: repository,
	}
}

func (a APIKeyAuthorizer) Authorize(r *http.Request) (*http.Request, error) {
	username, password, ok := r.BasicAuth()
	if !ok || len(username) == 0 {
		return nil, fmt.Errorf("wrong basic auth credentials provided")
	}

	user, err := a.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("cannot get user by username [%s]: %w", username, err)
	}

	if user.Username != username || password != mockPassword {
		return nil, fmt.Errorf("wrong username or password for username [%s]", user.Username)
	}

	return r.WithContext(context.WithValue(r.Context(), "user", user)), nil
}
