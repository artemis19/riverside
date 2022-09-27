package api

import (
	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func Serve() {
	// Creates gin router with default middleware
	Router = gin.Default()
	// Use CORS middleware since frontend could be separate from server location
	Router.Use(CORSMiddleware())

	Router.POST("/login", LoginHandler)
	Router.POST("/register", RegisterHandler)
	Router.Run(":8089")
}
