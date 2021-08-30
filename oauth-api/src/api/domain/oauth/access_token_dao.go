package oauth

import (
	"fmt"

	"github.com/Yapcheekian/rest-go/src/api/utils/errors"
)

var (
	tokens = make(map[string]*AccessToken)
)

func (at *AccessToken) Save() errors.ApiError {
	at.AccessToken = fmt.Sprintf("USR_%d", at.UserId)
	tokens[at.AccessToken] = at
	return nil
}

func GetAccessToken(accessToken string) (*AccessToken, errors.ApiError) {
	token := tokens[accessToken]
	if token == nil {
		return nil, errors.NewNotFoundApiError("access token not found")
	}
	return token, nil
}
