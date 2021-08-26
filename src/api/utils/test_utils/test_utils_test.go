package test_utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMockedContext(t *testing.T) {
	response := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "http://localhost:123/test", nil)

	assert.Nil(t, err)

	c := GetMockedContext(request, response)

	assert.EqualValues(t, http.MethodGet, c.Request.Method)
	assert.EqualValues(t, "123", c.Request.URL.Port())
	assert.EqualValues(t, "/test", c.Request.URL.Path)
}
