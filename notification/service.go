package notification

import (
	"context"
	"fmt"
	"test/dto"
)

type Worker struct {
	Ch  chan dto.Notification
	Ctx context.Context
}

func (w *Worker) Start() error {
	for {
		select {
		case notif := <-w.Ch:
			switch notif.Os {
			case "android":
				fmt.Println(notif)
				key := []string{notif.PushToken}
				err := SendPushNotification(key, notif.Title, notif.Body)
				if err != nil {
					return err
				}
			}
		case <-w.Ctx.Done():
			return nil
		}
	}
}
