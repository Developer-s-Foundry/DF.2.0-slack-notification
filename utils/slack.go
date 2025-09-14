package utils

import (
	"context"
	"time"

	"github.com/slack-go/slack"
)

func SendSlackNotification(api *slack.Client, message string) error {
	channelID := "C09FJSVR52M"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, _, err := api.PostMessageContext(
		ctx,
		channelID,
		slack.MsgOptionText(message, false),
	)
	if err != nil {
		return err
	}
	return nil
}
