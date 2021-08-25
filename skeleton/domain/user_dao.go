package domain

import (
	"fmt"
	"net/http"

	"github.com/Yapcheekian/rest-go/skeleton/utils"
)

var (
	users = map[int64]*User{
		123: {FirstName: "yap", LastName: "cheekian", Email: "cheekian@gmail.com"},
	}
	UserDao userDaoInterface
)

func init() {
	UserDao = &userDao{}
}

type userDao struct{}

type userDaoInterface interface {
	GetUser(int64) (*User, *utils.ApplicationError)
}

func (ud *userDao) GetUser(userId int64) (*User, *utils.ApplicationError) {
	if user := users[userId]; user != nil {
		return user, nil
	}
	return nil, &utils.ApplicationError{
		Message:    fmt.Sprintf("userId %d not found", userId),
		StatusCode: http.StatusNotFound,
		Code:       "not_found",
	}
}
