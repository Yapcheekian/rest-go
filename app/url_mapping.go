package app

import (
	"github.com/Yapcheekian/rest-go/controller"
)

func mapUrls() {
	router.GET("/users/:userId", controller.GetUser)
}
