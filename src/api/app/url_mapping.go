package app

import (
	"github.com/Yapcheekian/rest-go/src/api/controller/health"
	"github.com/Yapcheekian/rest-go/src/api/controller/repositories"
)

func mapUrls() {
	router.GET("/health", health.HealthCheck)
	router.POST("/repositories", repositories.CreateRepo)
}
