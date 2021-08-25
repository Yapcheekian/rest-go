package domain

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	t.Run("NoUserFound", func(t *testing.T) {
		user, err := UserDao.GetUser(0)

		assert.Nil(t, user)
		assert.NotNil(t, err)
		assert.EqualValues(t, http.StatusNotFound, err.StatusCode)
	})

	t.Run("NoError", func(t *testing.T) {
		user, err := UserDao.GetUser(123)

		assert.NotNil(t, user)
		assert.Nil(t, err)
		assert.EqualValues(t, "yap", user.FirstName)
	})
}
