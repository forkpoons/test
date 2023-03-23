package main

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/api/option"
	"log"
	"os"
	"time"
)

func main() {
	fmt.Println("qwe")
	link := "amqp://rabbitmq:rabbitmq@" + os.Getenv("NAME") + ":5672/"
	conn, err := amqp.Dial(link) // Создаем подключение к RabbitMQ
	if err != nil {
		log.Fatalf("unable to open connect to RabbitMQ server. Error: %s", err)
	}

	defer func() {
		_ = conn.Close() // Закрываем подключение в случае удачной попытки
	}()
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open channel. Error: %s", err)
	}

	defer func() {
		_ = ch.Close() // Закрываем канал в случае удачной попытки открытия
	}()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatalf("failed to declare a queue. Error: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Hello World!"
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		log.Fatalf("failed to publish a message. Error: %s", err)
	}

	log.Printf(" [x] Sent %s\n", body)

	ch2, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel. Error: %s", err)
	}

	defer func() {
		_ = ch2.Close() // Закрываем подключение в случае удачной попытки подключения
	}()

	q2, err := ch2.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatalf("failed to declare a queue. Error: %s", err)
	}

	messages, err := ch2.Consume(
		q2.Name, // queue
		"",      // consumer
		true,    // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		log.Fatalf("failed to register a consumer. Error: %s", err)
	}

	var forever chan struct{}

	go func() {
		for message := range messages {
			key := []string{"d2SZnL-nR7WrkVB-u4z5b9:APA91bGaOhx8sCJi2rq5GjfKhCIELIsnH97K12R4Q7F-71aXU1sdyQP_04s1DzMWy7Gye2hhJMx8iYKNHs4kjF2A5snvXQEwMfHtgBykkVmeYPwxfCAY6MPVIoj8kejjRcoTa0ugkTj-"}
			err := SendPushNotification(key)
			if err != nil {
				log.Printf("received a message: %s", err)
			}
			log.Printf("received a message: %s", message.Body)

		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func SendPushNotification(deviceTokens []string) error {
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
			Title: "Congratulations!!",
			Body:  "You have just implement push notification",
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
