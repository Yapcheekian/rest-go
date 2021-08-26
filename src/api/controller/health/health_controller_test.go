package health

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Yapcheekian/rest-go/src/api/utils/test_utils"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/health", nil)
	c := test_utils.GetMockedContext(request, response)

	HealthCheck(c)

	assert.EqualValues(t, http.StatusOK, response.Code)
	assert.EqualValues(t, "\"ok\"", response.Body.String())
}
