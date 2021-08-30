package oauth

import (
	"net/http"

	"github.com/Yapcheekian/rest-go/oauth-api/src/api/domain/oauth"
	"github.com/Yapcheekian/rest-go/oauth-api/src/api/service"
	"github.com/Yapcheekian/rest-go/src/api/utils/errors"
	"github.com/gin-gonic/gin"
)

func CreateAccessToken(c *gin.Context) {
	var request oauth.AccessTokenRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		apiErr := errors.NewBadRequestError("invalid json body")
		c.JSON(apiErr.GetStatus(), apiErr)
		return
	}

	token, err := service.OauthService.CreateAccessToken(request)

	if err != nil {
		c.JSON(err.GetStatus(), err)
		return
	}

	c.JSON(http.StatusCreated, token)
}

func GetAccessToken(c *gin.Context) {
	token, err := service.OauthService.GetAccessToken(c.Param("token_id"))

	if err != nil {
		c.JSON(err.GetStatus(), err)
		return
	}

	c.JSON(http.StatusOK, token)
}
