package services

import (
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/Yapcheekian/rest-go/src/api/clients/restclient"
	"github.com/Yapcheekian/rest-go/src/api/domain/repositories"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	restclient.StartMockUps()
	os.Exit(m.Run())
}

func TestCreateRepo(t *testing.T) {
	t.Run("invalid repository name", func(t *testing.T) {
		request := repositories.CreateRepoRequest{}
		resp, err := RepositoryService.CreateRepo(request)

		assert.Nil(t, resp)
		assert.NotNil(t, err)
		assert.EqualValues(t, "invalid repository name", err.GetMessage())
		assert.EqualValues(t, http.StatusBadRequest, err.GetStatus())
	})

	t.Run("github error", func(t *testing.T) {
		restclient.FlushMockUps()
		restclient.AddMockUp(restclient.Mock{
			Url:        "https://api.github.com/user/repos",
			HttpMethod: http.MethodPost,
			Response: &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       io.NopCloser(strings.NewReader(`{"message": "Requires authentication"}`)),
			},
		})

		request := repositories.CreateRepoRequest{
			Name: "test",
		}
		resp, err := RepositoryService.CreateRepo(request)

		assert.Nil(t, resp)
		assert.NotNil(t, err)
		assert.EqualValues(t, http.StatusUnauthorized, err.GetStatus())
	})

	t.Run("no error", func(t *testing.T) {
		restclient.FlushMockUps()
		restclient.AddMockUp(restclient.Mock{
			Url:        "https://api.github.com/user/repos",
			HttpMethod: http.MethodPost,
			Response: &http.Response{
				StatusCode: http.StatusCreated,
				Body:       io.NopCloser(strings.NewReader(`{"id": 123}`)),
			},
		})

		request := repositories.CreateRepoRequest{
			Name: "test",
		}
		resp, err := RepositoryService.CreateRepo(request)

		assert.NotNil(t, resp)
		assert.Nil(t, err)
		assert.EqualValues(t, 123, resp.Id)
	})
}
