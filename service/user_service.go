package service

import (
	"github.com/Yapcheekian/rest-go/domain"
	"github.com/Yapcheekian/rest-go/utils"
)

var (
	UserService userService
)

type userService struct{}

func (us *userService) GetUser(userId int64) (*domain.User, *utils.ApplicationError) {
	return domain.UserDao.GetUser(userId)
}
