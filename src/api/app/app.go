package app

import (
	"github.com/Yapcheekian/rest-go/src/api/log"
	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
)

func init() {
	router = gin.Default()
}

func StartApp() {
	log.Info("mapping urls...", "status:pending")
	mapUrls()
	log.Log.Info("url successfully mapped", "status:done")
	router.Run(":8080")
}
