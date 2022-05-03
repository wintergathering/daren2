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

	r.Static("/home", "./")
	r.LoadHTMLGlob("templates/*.html")

	//add a new dare
	r.POST("/home", func(c *gin.Context) {
		err := dareController.Save(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.HTML(http.StatusOK, "return.html", nil)
		}
	})

	//view all dares
	r.GET("/all_dares", dareController.ShowAll)

	//show a random dare
	r.GET("/rand_dare", dareController.ShowRandom)

	r.Run("localhost:8080")
}
