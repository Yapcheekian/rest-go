package githubprovider

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/Yapcheekian/rest-go/src/clients/restclient"
	"github.com/Yapcheekian/rest-go/src/domain/github"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	restclient.StartMockUps()
	os.Exit(m.Run())
}
func TestGetAuthorizationHeader(t *testing.T) {
	header := getAuthorizationHeader("ABC123")
	assert.EqualValues(t, "token ABC123", header)
}

func TestCreateRepo(t *testing.T) {
	t.Run("restclient send error", func(t *testing.T) {
		restclient.FlushMockUps()
		restclient.AddMockUp(restclient.Mock{
			Url:        "https://api.github.com/user/repos",
			HttpMethod: http.MethodPost,
			Err:        errors.New("invalid restclient response"),
		})

		response, err := CreateRepo("", github.CreateRepoRequest{})
		assert.Nil(t, response)
		assert.NotNil(t, err)
	})

	t.Run("invalid response body", func(t *testing.T) {
		restclient.FlushMockUps()
		invalidCloser, _ := os.Open("abc")

		restclient.AddMockUp(restclient.Mock{
			Url:        "https://api.github.com/user/repos",
			HttpMethod: http.MethodPost,
			Response: &http.Response{
				StatusCode: http.StatusCreated,
				Body:       invalidCloser,
			},
		})

		response, err := CreateRepo("", github.CreateRepoRequest{})

		assert.Nil(t, response)
		assert.NotNil(t, err)
	})

	t.Run("invalid response body with wrong json format", func(t *testing.T) {
		restclient.FlushMockUps()
		restclient.AddMockUp(restclient.Mock{
			Url:        "https://api.github.com/user/repos",
			HttpMethod: http.MethodPost,
			Response: &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       ioutil.NopCloser(strings.NewReader(`{"message": 1}`)),
			},
		})

		response, err := CreateRepo("", github.CreateRepoRequest{})

		assert.Nil(t, response)
		assert.NotNil(t, err)
	})
}
