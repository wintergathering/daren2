package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wintergathering/daren2/models"
	"github.com/wintergathering/daren2/repository"
)

type DareController interface {
	Save(c *gin.Context) error
	FindAll() []models.Dare
}

type controller struct {
	dare repository.DareRepository
}

func New(d repository.DareRepository) DareController {
	return &controller{
		dare: d,
	}
}

func (cn *controller) FindAll() {
	return cn.dare.FindAll()
	//RESUME HERE -- THIS NEEDS TO ALSO RETURN AN ERROR I THINK
}
