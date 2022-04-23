package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wintergathering/daren2/controller"
	"github.com/wintergathering/daren2/repository"
)

var (
	dareRepository repository.DareRepository = repository.NewDareRepository()
	dareController controller.DareController = controller.New(dareRepository)
)

func main() {

	r := gin.Default()

	r.Static("/", "./")
	r.POST("/", func(c *gin.Context) {
		err := dareController.Save(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "Dare is valid"})
		}
	})

	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })

	r.Run("localhost:8080")
}

// see this video for some help: https://www.youtube.com/watch?v=RHa4D6aNVpg
