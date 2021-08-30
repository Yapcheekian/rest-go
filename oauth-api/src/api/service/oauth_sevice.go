package service

import (
	"time"

	"github.com/Yapcheekian/rest-go/oauth-api/src/api/domain/oauth"
	"github.com/Yapcheekian/rest-go/src/api/utils/errors"
)

type oauthService struct{}

type oauthServiceInterface interface {
	CreateAccessToken(request oauth.AccessTokenRequest) (*oauth.AccessToken, errors.ApiError)
	GetAccessToken(accessToken string) (*oauth.AccessToken, errors.ApiError)
}

var (
	OauthService oauthServiceInterface
)

func init() {
	OauthService = &oauthService{}
}

func (o *oauthService) CreateAccessToken(request oauth.AccessTokenRequest) (*oauth.AccessToken, errors.ApiError) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	user, err := oauth.GetUserByUsernameAndPassword(request.Username, request.Password)

	if err != nil {
		return nil, err
	}

	token := oauth.AccessToken{
		UserId:  user.Id,
		Expires: time.Now().UTC().Add(24 * time.Hour).Unix(),
	}

	if err := token.Save(); err != nil {
		return nil, err
	}

	return &token, nil
}

func (o *oauthService) GetAccessToken(accessToken string) (*oauth.AccessToken, errors.ApiError) {
	token, err := oauth.GetAccessToken(accessToken)

	if err != nil {
		return nil, err
	}

	if token.IsExpire() {
		return nil, errors.NewNotFoundApiError("access token not found")
	}

	return token, nil
}
