package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run("localhost:8080")
}

// see this video for some help: https://www.youtube.com/watch?v=RHa4D6aNVpg
