package repositories

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Yapcheekian/rest-go/src/api/domain/repositories"
	"github.com/Yapcheekian/rest-go/src/api/services"
	"github.com/Yapcheekian/rest-go/src/api/utils/errors"
)

func CreateRepo(c *gin.Context) {
	var request repositories.CreateRepoRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		apiErr := errors.NewBadRequestError("invalid json body")
		c.JSON(apiErr.GetStatus(), apiErr)
		return
	}

	result, err := services.RepositoryService.CreateRepo(request)

	if err != nil {
		c.JSON(err.GetStatus(), err)
		return
	}
	c.JSON(http.StatusCreated, result)
}
