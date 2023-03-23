package rabbitmq

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Worker struct {
	Link string
	ch   chan Notification
	ctx  context.Context
}

type Notification struct {
	PushToken string
	Header    string
	Body      string
	Os        string
}

func (w *Worker) Start() error {

	conn, err := amqp.Dial(w.Link) // Создаем подключение к RabbitMQ
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel. Error: %s", err)
	}

	defer func() {
		_ = ch.Close() // Закрываем подключение в случае удачной попытки подключения
	}()

	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to declare a queue. Error: %s", err)
	}

	messages, err := ch.Consume(
		q.Name, // queue
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to register a consumer. Error: %s", err)
	}

	var forever chan struct{}

	go func() {
		for message := range messages {
			w.ch <- Notification{
				PushToken: "",
				Header:    "",
				Body:      "",
				Os:        "",
			}
			log.Printf("received a message: %s", message.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	return nil
}
