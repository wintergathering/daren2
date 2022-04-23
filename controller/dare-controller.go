package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wintergathering/daren2/models"
	"github.com/wintergathering/daren2/repository"
)

type DareController interface {
	Save(c *gin.Context) error
	FindAll() ([]models.Dare, error)
}

type controller struct {
	dare repository.DareRepository
}

func New(d repository.DareRepository) DareController {
	return &controller{
		dare: d,
	}
}

func (cn *controller) FindAll() ([]models.Dare, error) {
	return cn.dare.FindAll()
}

func (cn *controller) Save(c *gin.Context) error {

	var newDare *models.Dare

	dareTitle := c.PostForm("title")
	dareText := c.PostForm("text")

	newDare = &models.Dare{dareTitle, dareText}

	_, err := cn.dare.Save(newDare)

	return err
}

func (cn *controller) ShowAll(c *gin.Context) {
	dares, err := cn.dare.FindAll()

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "dares not retrieved"})
		return
	}

	data := gin.H{
		"title": "All of the Dares",
		"dares": dares,
	}

	c.HTML(http.StatusOK, "all_dares.html", data)
}

//resume here. next step is likely to make the html templates to try this out
