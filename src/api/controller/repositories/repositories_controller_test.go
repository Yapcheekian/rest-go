package repositories

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Yapcheekian/rest-go/src/api/clients/restclient"
	"github.com/Yapcheekian/rest-go/src/api/domain/repositories"
	"github.com/Yapcheekian/rest-go/src/api/utils/errors"
	"github.com/Yapcheekian/rest-go/src/api/utils/test_utils"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	restclient.StartMockUps()
	os.Exit(m.Run())
}

func TestCreateRepo(t *testing.T) {
	t.Run("invalid json request", func(t *testing.T) {
		response := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodPost, "/repositories", strings.NewReader(""))
		c := test_utils.GetMockedContext(request, response)
		CreateRepo(c)

		assert.EqualValues(t, http.StatusBadRequest, response.Code)

		apiErr, err := errors.NewApiErrorFromBytes(response.Body.Bytes())

		assert.NotNil(t, apiErr)
		assert.Nil(t, err)
		assert.EqualValues(t, http.StatusBadRequest, apiErr.GetStatus())
		assert.EqualValues(t, "invalid json body", apiErr.GetMessage())
	})

	t.Run("error from github", func(t *testing.T) {
		restclient.FlushMockUps()

		response := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodPost, "/repositories", strings.NewReader(`{"name": "test"}`))
		c := test_utils.GetMockedContext(request, response)

		restclient.AddMockUp(restclient.Mock{
			Url:        "https://api.github.com/user/repos",
			HttpMethod: http.MethodPost,
			Response: &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       io.NopCloser(strings.NewReader(`{"message": "Requires authentication"}`)),
			},
		})

		CreateRepo(c)

		assert.EqualValues(t, http.StatusUnauthorized, response.Code)

		apiErr, err := errors.NewApiErrorFromBytes(response.Body.Bytes())

		assert.EqualValues(t, http.StatusUnauthorized, apiErr.GetStatus())
		assert.Nil(t, err)
		assert.EqualValues(t, "Requires authentication", apiErr.GetMessage())
	})

	t.Run("no error", func(t *testing.T) {
		restclient.FlushMockUps()

		response := httptest.NewRecorder()
		request, _ := http.NewRequest(http.MethodPost, "/repositories", strings.NewReader(`{"name": "test"}`))
		c := test_utils.GetMockedContext(request, response)

		restclient.AddMockUp(restclient.Mock{
			Url:        "https://api.github.com/user/repos",
			HttpMethod: http.MethodPost,
			Response: &http.Response{
				StatusCode: http.StatusCreated,
				Body:       io.NopCloser(strings.NewReader(`{"id": 123, "name": "john"}`)),
			},
		})

		CreateRepo(c)

		assert.EqualValues(t, http.StatusCreated, response.Code)

		var result repositories.CreateRepoResponse

		err := json.Unmarshal(response.Body.Bytes(), &result)

		assert.Nil(t, err)
		assert.EqualValues(t, 123, result.Id)
		assert.EqualValues(t, "john", result.Name)
		assert.EqualValues(t, "", result.Owner)
	})
}
