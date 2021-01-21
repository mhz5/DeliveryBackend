package location

import (
	"context"
	"fmt"
	"log"
	"math"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"firebase.google.com/go/auth"
	"github.com/pkg/errors"
	"github.com/uber/h3-go/v3"

	lpb "fda/proto/location"
	"fda/server/gcloud"
)


const driverLocationTopic = "driver-location"
const driverLocationsCollectionName = "driverLocations"

var (
	pubsubClient    *pubsub.Client
	firestoreClient *firestore.Client
	authClient      *auth.Client
)

func init() {
	ctx := context.Background()
	pubsubClient = gcloud.GetPubsubClient(ctx)
	firestoreClient = gcloud.GetFirestoreClient(ctx)
	authClient = gcloud.GetAuthClient(ctx)
}

type Server struct{}

// TODO: Implement a pubsub proxy on GKE.
func (s *Server) PublishLocation(ctx context.Context, in *lpb.PublishLocationRequest) (*lpb.PublishLocationResponse, error) {
	log.Printf("Handling PublishLocation request [%v] with context %v", in, ctx)

	err := publishToPubSub(ctx, in)
	if err != nil {
		errors.Wrap(err, "failed to publish location to pub sub")
	}

	token, err := authClient.VerifyIDToken(ctx, in.UserToken)
	if err != nil {
		return &lpb.PublishLocationResponse{}, errors.Wrap(err, "could not verify ID token")
	}
	err = updateDriverLocation(ctx, token.UID, in.Location.Lat, in.Location.Long)
	if err != nil {
		errors.Wrap(err, "failed to update driver location")
	}

	return &lpb.PublishLocationResponse{Error: ""}, err
}

func publishToPubSub(ctx context.Context, in *lpb.PublishLocationRequest) error {
	topic := pubsubClient.Topic(driverLocationTopic)
	res := topic.Publish(ctx, &pubsub.Message{
		Data: []byte(fmt.Sprintf(
			"location (%f, %f)",
			in.GetLocation().GetLat(), in.GetLocation().GetLong())),
	})
	_, err := res.Get(ctx)
	return err
}

type location struct {
	UserId string `firestore:"userId"`
	H3Id   int64  `firestore:"h3IdRes11"` // Remember to change this tag if the resolution changes.
}

func updateDriverLocation(ctx context.Context, userId string, lat float64, lng float64) error {
	collection := firestoreClient.Collection(driverLocationsCollectionName)
	_, err := collection.Doc(userId).Set(ctx,
		location{
			H3Id: int64(h3FromLatLng(lat, lng)),
		})

	return errors.Wrap(err, "could not update driver location")
}

const Resolution25Meters = 11

func test() {
	lat := 38.628124411058806
	lng := -90.32835026159316
	h3Play(lat, lng)
}

func h3FromLatLng(lat float64, lng float64) h3.H3Index {
	geo := h3.GeoCoord{
		Latitude: lat,
		Longitude: lng,
	}
	// Resolutions: https://h3geo.org/docs/core-library/restable
	return h3.FromGeo(geo, Resolution25Meters)
}

func h3Play(lat float64, lng float64) (h3.H3Index, h3.H3Index) {
	center := h3FromLatLng(lat, lng)
	ring := h3.KRing(center, 10)
	return minMax(ring)
}

func minMax(arr []h3.H3Index) (h3.H3Index, h3.H3Index) {
	min := h3.H3Index(math.MaxInt64)
	max := h3.H3Index(0)
	for _, v := range arr {
		if v > max { max = v }
		if v < min { min = v }
	}
	return min, max
}

