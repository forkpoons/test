package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"test/dto"
)

type Worker struct {
	Link string
	Ch   chan dto.Notification
	Ctx  context.Context
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
			n := dto.Notification{}
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
