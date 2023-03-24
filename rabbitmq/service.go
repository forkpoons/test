package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Worker struct {
	Link string
	Ch   chan Notification
	Ctx  context.Context
}

type Notification struct {
	PushToken string `json:"push_token"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Os        string `json:"os"`
}

func (w *Worker) Start() error {
	fmt.Println("rabbit")
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
	for {
		select {
		case message := <-messages:
			fmt.Println(message.Body)
			n := Notification{}
			err := json.Unmarshal(message.Body, &n)
			if err != nil {
				return err
			}
			w.Ch <- n
			log.Printf("received a message: %s", message.Body)
		case <-w.Ctx.Done():
			return nil
		}
	}
	return nil
}
