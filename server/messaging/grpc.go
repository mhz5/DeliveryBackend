package messaging

import (
	"context"
	"fmt"

	"google.golang.org/api/fcm/v1"
)

var client *fcm.Service

func init() {
	ctx := context.Background()
	var err error
	client, err = fcm.NewService(ctx)
	if err != nil {
		panic(fmt.Errorf("could not create Firebase Cloud Messaging client: %v", err))
	}
}

func sendPushNotification(fcmToken string) {
}
