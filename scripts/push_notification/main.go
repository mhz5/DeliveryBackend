package main

import (
	"context"

	"log"
	"firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

const mikeAndroidToken="etvClaI2QeiE-w9hWa75CK:APA91bFgcvs5ngH68ve84wmCCQzDBdOrj69Ac1qrlLIVOS7fSEBCIBuaYX1Vr49ZN6ZIckg9xhFpy4xsasrehUUk16RS3rJlfvNDqPbgBCmHooVe8l8W8JGl9HLba9J1xcOELxOfCgqF"
const mikeIosToken="eDmGmlUMS0GpuVhPOtBJr2:APA91bEklggxumpGAorAQmkgwvpRW-C414T6CAEd0We6g_9V2hD_4WSszZMGVskFQVyj1zRUabP35zcKrdOyUygfstM9_YsjKw5MIxZlxe78ClNRzqy7ReqzbZ2SPm4a2t47FGijQZj9"
const ANDROID=0
const IOS=1

const PLATFORM=IOS

func main() {
	config := &firebase.Config{
		ProjectID: "food-prod",
	}
	ctx := context.Background()
	firebaseApp, err := firebase.NewApp(ctx, config)
	if err != nil {
		log.Fatal("Cannot create new Firebase app")
	}
	firebaseMessaging, err := firebaseApp.Messaging(ctx)
	if err != nil {
		log.Fatal("Cannot access Firebase Messaging")
	}

	msg := createMessage(PLATFORM)
	firebaseMessaging.Send(ctx, msg)
}

func createMessage(phoneType int) *messaging.Message {
	if phoneType == ANDROID {
		androidConfig := &messaging.AndroidConfig{}
		notification := &messaging.AndroidNotification{}
		notification.Body = "Test Body"
		notification.ChannelID = "ring-sound"
		androidConfig.Notification = notification

		return &messaging.Message{
			Token: mikeAndroidToken,
			Android: androidConfig,
		}
	}

	return &messaging.Message{
		Token:   mikeIosToken,
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{
				"apns-push-type": "background",
				"apns-priority":  "5",
				"apns-topic":     "org.reactjs.native.example.FdaDriverApp",
			},
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{ContentAvailable: true},
			},
		},
	}
}
