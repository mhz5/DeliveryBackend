package gcloud

import (
	"context"
	"sync"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/pkg/errors"
)

const (
	projectId = "food-prod"
)

var (
	pubsubOnce      sync.Once
	firestoreOnce   sync.Once
	firebaseOnce    sync.Once
	authOnce	    sync.Once

	pubsubClient    *pubsub.Client
	firestoreClient *firestore.Client
	firebaseApp 	*firebase.App
	authClient      *auth.Client
)

func init() {
	ctx := context.Background()
	GetPubsubClient(ctx)
	GetFirestoreClient(ctx)
	GetAuthClient(ctx)
}

// TODO: Is there a better way than panic to handle errors?
// GetPubsubClient returns the singleton instance of the google cloud pubsub client.
// On the first ever call, it creates the client instance.
// Thread safe.
func GetPubsubClient(ctx context.Context) *pubsub.Client {
	pubsubOnce.Do(func() {
		pubsubClient = newPubsubClient(ctx)
	})
	return pubsubClient
}

// GetFirestoreClient returns the singleton instance of the google cloud firestore client.
// On the first ever call, it creates the client instance.
// Thread safe.
func GetFirestoreClient(ctx context.Context) *firestore.Client {
	firestoreOnce.Do(func() {
		firestoreClient = newFirestoreClient(ctx)
	})
	return firestoreClient
}

// GetAuthClient returns the singleton instance of the google cloud firebase auth client.
// On the first ever call, it creates the client instance.
// Thread safe.
func GetAuthClient(ctx context.Context) *auth.Client {
	authOnce.Do(func() {
		authClient = newAuthClient(ctx)
	})
	return authClient
}

func newPubsubClient(ctx context.Context) *pubsub.Client {
	client, err := pubsub.NewClient(ctx, "food-prod")
	if err != nil {
		panic(errors.Wrap(err, "could not create pubsub client"))
	}
	return client
}

func newFirestoreClient(ctx context.Context) *firestore.Client {
	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		panic(errors.Wrap(err, "could not create firestore client"))
	}
	return client
}

func newAuthClient(ctx context.Context) *auth.Client {
	var err error
	firebaseOnce.Do(func() {
		firebaseApp = newFirebaseApp(ctx)
	})
	client, err := firebaseApp.Auth(ctx)
	if err != nil {
		panic(errors.Wrap(err, "cannot create firebase auth client"))
	}
	return client
}

func newFirebaseApp(ctx context.Context) *firebase.App {
	config := &firebase.Config{
		ProjectID: projectId,
	}
	app, err := firebase.NewApp(ctx, config)
	if err != nil {
		panic(errors.Wrap(err, "cannot create new firebase app"))
	}
	return app
}
