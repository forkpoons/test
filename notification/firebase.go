package notification

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	"google.golang.org/api/option"
)

func SendPushNotification(deviceTokens []string, title string, body string) error {
	//decodedKey, err := getDecodedFireBaseKey()
	//fmt.Println("ww")
	//if err != nil {
	//	return err
	//}

	opt := option.WithCredentialsFile("service-account-other.json")
	//opts := []option.ClientOption{option.WithCredentialsJSON(decodedKey)}

	app, err := firebase.NewApp(context.Background(), nil, opt)

	if err != nil {
		fmt.Printf("Error in initializing firebase : %s\n", err)
		return err
	}

	fcmClient, err := app.Messaging(context.Background())

	if err != nil {
		return err
	}

	response, err := fcmClient.SendMulticast(context.Background(), &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Tokens: deviceTokens,
	})

	if err != nil {
		return err
	}

	fmt.Println("Response success count : ", response.SuccessCount)
	fmt.Println("Response failure count : ", response.FailureCount)

	return nil
}
