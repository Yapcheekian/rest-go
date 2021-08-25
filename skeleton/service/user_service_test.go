package service

import (
	"net/http"
	"testing"

	"github.com/Yapcheekian/rest-go/skeleton/domain"
	"github.com/Yapcheekian/rest-go/skeleton/utils"
	"github.com/stretchr/testify/assert"
)

var (
	getUserFunction func(int64) (*domain.User, *utils.ApplicationError)
)

type usersDaoMock struct{}

func init() {
	domain.UserDao = &usersDaoMock{}
}

func (m *usersDaoMock) GetUser(userId int64) (*domain.User, *utils.ApplicationError) {
	return getUserFunction(userId)
}

func TestGetUser(t *testing.T) {
	t.Run("UserNotFoundInDatabase", func(t *testing.T) {
		getUserFunction = func(userId int64) (*domain.User, *utils.ApplicationError) {
			return nil, &utils.ApplicationError{
				StatusCode: http.StatusNotFound,
				Message:    "User 0 does not exists",
			}
		}

		user, err := UserService.GetUser(0)

		assert.Nil(t, user)
		assert.NotNil(t, err)
	})

	t.Run("NoError", func(t *testing.T) {
		getUserFunction = func(userId int64) (*domain.User, *utils.ApplicationError) {
			return &domain.User{
				FirstName: "test",
				LastName:  "test",
				Email:     "test",
			}, nil
		}

		user, err := UserService.GetUser(123)

		assert.NotNil(t, user)
		assert.Nil(t, err)
		assert.EqualValues(t, "test", user.FirstName)
		assert.EqualValues(t, "test", user.LastName)
		assert.EqualValues(t, "test", user.Email)
	})
}
