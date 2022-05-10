package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/wintergathering/daren2/controller"
	"github.com/wintergathering/daren2/repository"
)

var (
	dareRepository repository.DareRepository = repository.NewDareRepository()
	dareController controller.DareController = controller.New(dareRepository)
)

func main() {

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

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

	r.Run(":" + httpPort)
}
