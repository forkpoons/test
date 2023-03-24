package main

import (
	"context"
	"os"
	"test/dto"
	"test/notification"
	"test/rabbitmq"
)

func main() {
	ctx := context.Background()
	rabbitChannel := make(chan dto.Notification)
	rabbit := rabbitmq.Worker{
		Link: "amqp://" + os.Getenv("LOGIN") + ":" + os.Getenv("PASS") + "@" + os.Getenv("LINK") + ":" + os.Getenv("PORT") + "/",
		Ch:   rabbitChannel,
		Ctx:  ctx,
	}
	go func() {
		err := rabbit.Start()
		if err != nil {
			return
		}
	}()
	notif := notification.Worker{
		Ch:  rabbitChannel,
		Ctx: ctx,
	}
	go func() {
		err := notif.Start()
		if err != nil {
			return
		}
	}()

}
