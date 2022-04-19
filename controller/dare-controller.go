package controller

import (
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

	dareTitle := c.PostForm("title")
	dareText := c.PostForm("text")

	newDare := &models.Dare{dareTitle, dareText}

	cn.dare.Save(newDare)
	//resume here later
}
