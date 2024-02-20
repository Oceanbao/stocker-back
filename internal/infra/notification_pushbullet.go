package infra

import (
	"context"
	"errors"
	"os"

	"github.com/joho/godotenv"
	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/pushbullet"
)

type NotifierPushbullet struct {
	service notify.Notifier
}

func NewNotifierPushbullet() (*NotifierPushbullet, error) {
	_ = godotenv.Load()
	accessToken, ok := os.LookupEnv("PUSHBULLET_TOKEN")
	if !ok {
		return nil, errors.New("missing pushbullet access token env")
	}

	notifier := notify.New()
	service := pushbullet.New(accessToken)
	service.AddReceivers("Chrome")
	notifier.UseServices(service)

	return &NotifierPushbullet{
		service: notifier,
	}, nil
}

func (n *NotifierPushbullet) Sendf(topic, msg string) {
	_ = n.service.Send(context.Background(), topic, msg)
}
