package app

import (
	"github.com/Yapcheekian/rest-go/oauth-api/src/api/controller/oauth"
	"github.com/Yapcheekian/rest-go/src/api/controller/health"
)

func mapUrls() {
	router.GET("/health", health.HealthCheck)
	router.POST("/oauth/access_token", oauth.CreateAccessToken)
	router.GET("/oauth/access_token/:token_id", oauth.GetAccessToken)
}
