package service

import "github.com/wintergathering/daren2/models"

type DareService interface {
	Save(models.Dare) models.Dare
}

type dareService struct {
	dares []models.Dare
}

func New() DareService {
	return &dareService{}
}

func (ds *dareService) Save(d models.Dare) models.Dare {
	//add save func here -- save to firestore
}
