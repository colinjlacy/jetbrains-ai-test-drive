package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Start() {
	r := gin.Default()

	registerHealthEndpoint(r)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello, world",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}

func registerHealthEndpoint(engine *gin.Engine) {
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
}
