package repository

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/wintergathering/daren2/models"
	"google.golang.org/api/iterator"
)

type DareRepository interface {
	Save(d *models.Dare) (*models.Dare, error)
	FindAll() ([]models.Dare, error)
}

type repo struct{}

const (
	projectID      string = "dares-app-346910"
	collectionName string = "dares"
)

func NewDareRepository() DareRepository {
	return &repo{}
}

func (*repo) Save(d *models.Dare) (*models.Dare, error) {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, projectID)

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return nil, err
	}

	defer client.Close()

	_, _, err = client.Collection(collectionName).Add(ctx, map[string]interface{}{
		"Title": d.Title,
		"Text":  d.Text,
	})

	if err != nil {
		log.Fatalf("Failed adding a new dare: %v", err)
		return nil, err
	}

	return d, nil
}

func (*repo) FindAll() ([]models.Dare, error) {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, projectID)

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return nil, err
	}

	defer client.Close()

	var dares []models.Dare

	iter := client.Collection(collectionName).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		dare := models.Dare{
			Title: doc.Data()["Title"].(string),
			Text:  doc.Data()["Text"].(string),
		}

		dares = append(dares, dare)
	}

	return dares, nil
}
