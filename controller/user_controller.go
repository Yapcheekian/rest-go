package controller

import (
	"net/http"
	"strconv"

	"github.com/Yapcheekian/rest-go/service"
	"github.com/Yapcheekian/rest-go/utils"
	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Param("userId"), 10, 64)
	if err != nil {
		apiErr := &utils.ApplicationError{
			Message:    "userId must be a number",
			StatusCode: http.StatusBadRequest,
			Code:       "bad_request",
		}
		utils.RespondError(c, apiErr)
		return
	}

	user, apiErr := service.UserService.GetUser(userId)

	if apiErr != nil {
		utils.RespondError(c, apiErr)
		return
	}

	utils.Respond(c, http.StatusOK, user)
}
