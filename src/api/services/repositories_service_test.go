package services

import (
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
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

func TestCreateRepoConcurrent(t *testing.T) {
	t.Run("invalid request", func(t *testing.T) {
		request := repositories.CreateRepoRequest{}
		output := make(chan repositories.CreateRepositoryResult)

		service := &repoService{}
		go service.createRepoConcurrent(request, output)

		result := <-output

		assert.NotNil(t, result)
		assert.Nil(t, result.Response)
		assert.EqualValues(t, "invalid repository name", result.Error.GetMessage())
		assert.EqualValues(t, http.StatusBadRequest, result.Error.GetStatus())
	})

	t.Run("error from github", func(t *testing.T) {
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
		output := make(chan repositories.CreateRepositoryResult)

		service := &repoService{}
		go service.createRepoConcurrent(request, output)

		result := <-output

		assert.NotNil(t, result)
		assert.Nil(t, result.Response)
		assert.EqualValues(t, "Requires authentication", result.Error.GetMessage())
		assert.EqualValues(t, http.StatusUnauthorized, result.Error.GetStatus())
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
		output := make(chan repositories.CreateRepositoryResult)

		service := &repoService{}
		go service.createRepoConcurrent(request, output)

		result := <-output

		assert.NotNil(t, result)
		assert.NotNil(t, result.Response)
		assert.EqualValues(t, 123, result.Response.Id)
	})
}

func TestHandleRepoResults(t *testing.T) {
	input := make(chan repositories.CreateRepositoryResult)
	output := make(chan repositories.CreateReposResponse)
	var wg sync.WaitGroup

	service := &repoService{}

	go service.handleRepoResults(&wg, input, output)

	wg.Add(1)
	go func() {
		input <- repositories.CreateRepositoryResult{
			Response: &repositories.CreateRepoResponse{
				Id:    123,
				Name:  "test",
				Owner: "John",
			},
		}
	}()

	wg.Wait()
	close(input)

	result := <-output

	assert.NotNil(t, result)
	assert.EqualValues(t, "test", result.Results[0].Response.Name)
	assert.EqualValues(t, "John", result.Results[0].Response.Owner)
	assert.EqualValues(t, 123, result.Results[0].Response.Id)
}

func TestCreateRepos(t *testing.T) {
	t.Run("invalid requests", func(t *testing.T) {
		requests := []repositories.CreateRepoRequest{
			{},
			{Name: " "},
		}
		response, err := RepositoryService.CreateRepos(requests)

		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.EqualValues(t, http.StatusBadRequest, response.StatusCode)
		assert.EqualValues(t, 2, len(response.Results))
		assert.EqualValues(t, "invalid repository name", response.Results[0].Error.GetMessage())
		assert.EqualValues(t, "invalid repository name", response.Results[1].Error.GetMessage())
	})

	t.Run("partially success request", func(t *testing.T) {
		restclient.FlushMockUps()
		restclient.AddMockUp(restclient.Mock{
			Url:        "https://api.github.com/user/repos",
			HttpMethod: http.MethodPost,
			Response: &http.Response{
				StatusCode: http.StatusCreated,
				Body:       io.NopCloser(strings.NewReader(`{"id": 123}`)),
			},
		})

		requests := []repositories.CreateRepoRequest{
			{},
			{Name: "test", Description: "this is a description"},
		}
		response, err := RepositoryService.CreateRepos(requests)

		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.EqualValues(t, http.StatusPartialContent, response.StatusCode)
		assert.EqualValues(t, 2, len(response.Results))

		for _, result := range response.Results {
			if result.Error != nil {
				assert.EqualValues(t, "invalid repository name", result.Error.GetMessage())
				assert.EqualValues(t, http.StatusBadRequest, result.Error.GetStatus())
				continue
			}
			assert.EqualValues(t, 123, result.Response.Id)
		}
	})

	t.Run("all success requests", func(t *testing.T) {
		restclient.FlushMockUps()
		restclient.AddMockUp(restclient.Mock{
			Url:        "https://api.github.com/user/repos",
			HttpMethod: http.MethodPost,
			Response: &http.Response{
				StatusCode: http.StatusCreated,
				Body:       io.NopCloser(strings.NewReader(`{"id": 123}`)),
			},
		})

		requests := []repositories.CreateRepoRequest{
			{Name: "test", Description: "this is a description"},
		}

		response, err := RepositoryService.CreateRepos(requests)

		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.EqualValues(t, http.StatusCreated, response.StatusCode)
		assert.EqualValues(t, 1, len(response.Results))
	})
}
