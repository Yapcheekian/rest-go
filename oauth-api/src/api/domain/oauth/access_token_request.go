package oauth

import (
	"strings"

	"github.com/Yapcheekian/rest-go/src/api/utils/errors"
)

type AccessTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *AccessTokenRequest) Validate() errors.ApiError {
	if strings.TrimSpace(r.Username) == "" {
		return errors.NewBadRequestError("invalid username")
	}
	if strings.TrimSpace(r.Password) == "" {
		return errors.NewBadRequestError("invalid password")
	}
	return nil
}
