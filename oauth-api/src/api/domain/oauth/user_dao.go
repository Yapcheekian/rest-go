package oauth

import (
	"fmt"

	"github.com/Yapcheekian/rest-go/src/api/utils/errors"
)

const (
	queryGetUserByUsernameAndPassword = "SELECT id, username FROM users where username = ? and password = ?;"
)

var (
	users = map[string]*User{
		"yap": {Id: 123, Username: "cheekian"},
	}
)

func GetUserByUsernameAndPassword(username string, password string) (*User, errors.ApiError) {
	user := users[username]
	if user == nil {
		return nil, errors.NewNotFoundApiError(fmt.Sprintf("username %s not found", username))
	}
	return user, nil
}
