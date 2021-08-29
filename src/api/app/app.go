package app

import (
	"github.com/Yapcheekian/rest-go/src/api/log/logrus"
	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
)

func init() {
	router = gin.Default()
}

func StartApp() {
	logrus.Info("mapping urls...", "status:pending")
	mapUrls()
	logrus.Log.Info("url successfully mapped", "status:done")
	router.Run(":8080")
}
