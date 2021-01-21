package restaurant

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"

	rpb "fda/proto/restaurant"
	"fda/server/gcloud"
)

var (
	firestoreClient *firestore.Client
)

func init() {
	ctx := context.Background()
	firestoreClient = gcloud.GetFirestoreClient(ctx)
}

type Server struct{}

func (s *Server) GetRestaurants(ctx context.Context, in *rpb.GetRestaurantsRequest) (*rpb.GetRestaurantsResponse, error) {
	iter := getRestaurantIterator(ctx)
	docs, err := iter.GetAll()
	if err != nil {
		return &rpb.GetRestaurantsResponse{}, errors.Wrap(err, "failed to query restaurants from firestore")
	}
	defer iter.Stop()

	restaurants := parseRestaurants(docs, ctx)

	return &rpb.GetRestaurantsResponse{Restaurants: restaurants}, nil
}

func getRestaurantIterator(ctx context.Context) (*firestore.DocumentIterator) {
	restaurants := firestoreClient.Collection("restaurants")
	// TODO: This logic is wrong. Currently just returns 5 restaurants.
	query := restaurants.Where("id", ">", 0)
	query = query.Where("id", "<", 6)
	iter := query.Documents(ctx)
	return iter
}

func parseRestaurants(docs []*firestore.DocumentSnapshot, ctx context.Context) ([]*rpb.Restaurant) {
	numRestaurants := len(docs)
	restaurants := make([]*rpb.Restaurant, numRestaurants)
	for i, doc := range docs {
		r, err := parseRestaurant(ctx, doc)
		if err != nil {
			log.Printf("failed to parse %d-th restaurant: %v", i, err)
		}
		restaurants[i] = r
	}
	return restaurants
}
