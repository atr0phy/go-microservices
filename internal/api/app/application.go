package app

import (
	"github.com/gin-gonic/gin"
)

var (
	router *gin.Engine
)

func SetupRouter() *gin.Engine {
	router = gin.Default()
	mapUrls()
	return router
}
