package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wintergathering/daren2/models"
	"github.com/wintergathering/daren2/repository"
)

type DareController interface {
	Save(c *gin.Context) error
	FindAll() ([]models.Dare, error)
	ShowAll(c *gin.Context)
	ShowRandom(c *gin.Context)
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
	dareSeen := false

	newDare = &models.Dare{dareTitle, dareText, dareSeen}

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
		"hd":    "All of the Dares",
		"dares": dares,
	}

	c.HTML(http.StatusOK, "all_dares.html", data)
}

func (cn *controller) ShowRandom(c *gin.Context) {
	d, id, err := cn.dare.GetRandDare()

	if err != nil {
		//c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		c.HTML(http.StatusOK, "no_dares.html", nil)
		return
	}

	c.HTML(http.StatusOK, "single_dare.html", d)

	//update the dare to be seen after it's pulled
	_, err = cn.dare.UpdateSeen(id)

	if err != nil {
		log.Fatalf("error updating seen value of dare: %v", err)
	}
}
