package frstr

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/wintergathering/daren2"
	"google.golang.org/api/iterator"
)

const (
	projectID = "dares-app-346910"
	daresColl = "dares"
)

type dareService struct {
	Client *firestore.Client
}

// function to create a new firestore client
func NewFirestoreClient() *firestore.Client {
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, projectID)

	if err != nil {
		log.Fatal("Couldn't create firestore client!")
	}

	return client
}

// function to create a new dareService that implements the DareService interface
func NewDareService(client *firestore.Client) daren2.DareService {
	return &dareService{
		Client: client,
	}
}

func (ds dareService) CreateDare(ctx context.Context, dare *daren2.Dare) error {
	_, _, err := ds.Client.Collection(daresColl).Add(ctx, dare)

	if err != nil {
		return err
	}

	return nil
}

func (ds dareService) GetRandomDare(ctx context.Context) (*daren2.Dare, error) {

	iter := ds.Client.Collection(daresColl).Where("Seen", "==", false).Limit(1).Documents(ctx)

	defer iter.Stop()

	doc, err := iter.Next()

	if err == iterator.Done {
		return nil, daren2.ErrNoDare
	}

	if err != nil {
		return nil, err
	}

	var dare *daren2.Dare

	doc.DataTo(&dare)

	//update dare as seen
	id := doc.Ref.ID

	updateSeen := []firestore.Update{
		{Path: "Seen", Value: true},
	}
	_, err = ds.Client.Collection(daresColl).Doc(id).Update(ctx, updateSeen)

	if err != nil {
		return nil, err
	}

	return dare, nil
}

func (ds dareService) GetAllDares(ctx context.Context) ([]*daren2.Dare, error) {
	docs, err := ds.Client.Collection(daresColl).Documents(ctx).GetAll()

	if err != nil {
		return nil, err
	}

	var dares []*daren2.Dare

	for _, doc := range docs {
		var d *daren2.Dare

		doc.DataTo(&d)

		dares = append(dares, d)
	}

	return dares, nil
}
