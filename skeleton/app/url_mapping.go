package app

import (
	"github.com/Yapcheekian/rest-go/skeleton/controller"
)

func mapUrls() {
	router.GET("/users/:userId", controller.GetUser)
}
