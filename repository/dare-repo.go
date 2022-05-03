package repository

import (
	"context"
	"log"
	"math/rand"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/wintergathering/daren2/models"
	"google.golang.org/api/iterator"
)

type DareRepository interface {
	Save(d *models.Dare) (*models.Dare, error)
	FindAll() ([]models.Dare, error)
	GetRandDare() (*models.Dare, string, error)
	UpdateSeen(id string) (string, error)
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
		"Seen":  d.Seen,
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
			Seen:  doc.Data()["Seen"].(bool),
		}

		dares = append(dares, dare)
	}

	return dares, nil
}

//find by name
func findByID(id string) (*models.Dare, error) {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, projectID)

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return nil, err
	}

	defer client.Close()

	dsnap, err := client.Collection(collectionName).Doc(id).Get(ctx)

	dare := &models.Dare{
		Title: dsnap.Data()["Title"].(string),
		Text:  dsnap.Data()["Text"].(string),
		Seen:  dsnap.Data()["Seen"].(bool),
	}

	return dare, nil

}

func getRandID() (string, error) {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, projectID)

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return "", err
	}

	defer client.Close()

	var ids []string

	iter := client.Collection(collectionName).Where("Seen", "==", false).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Error iterating over documents: %v", err)
		}
		idx := doc.Ref.ID

		ids = append(ids, idx)
	}

	rand.Seed(time.Now().Unix())

	n := rand.Intn(len(ids))

	id := ids[n]

	return id, nil
}

func (*repo) GetRandDare() (*models.Dare, string, error) {
	id, err := getRandID()

	if err != nil {
		log.Fatalf("Error getting a random id: %v", err)
	}

	dare, err := findByID(id)

	if err != nil {
		log.Fatalf("Error getting random dare: %v", err)
	}

	return dare, id, nil
}

func (*repo) UpdateSeen(id string) (string, error) {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, projectID)

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	_, err = client.Collection(collectionName).Doc(id).Update(ctx, []firestore.Update{
		{Path: "Seen", Value: true},
	})

	if err != nil {
		log.Fatalf("Failed to update `seen` field: %v", err)
	}

	msg := "Updated seen field"

	return msg, nil
}
