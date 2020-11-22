package main

import (
	"fmt"
	"os"

	"github.com/jhunt/go-chat"
)

func main() {
	var (
		bot chat.Bot
		err error
	)

	channels := []string{}
	if token := os.Getenv("BOT_SLACK_TOKEN"); token != "" {
		fmt.Printf("connecting to slack...\n")
		bot, err = chat.Slack(token)

		c := os.Getenv("BOT_SLACK_CHANNEL")
		if c == "" {
			c = "testing"
		}
		channels = append(channels, c)

	} else {
		fmt.Printf("connecting to tty...\n")
		bot, err = chat.Terminal()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect failed: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("connected; awaiting *every* message...\n")
	bot.Every(func(msg chat.Message) chat.Then {
		fmt.Printf("%s [%s] %s: %s\n", msg.Received.Format("2006-01-02 15:04:05+0000"), msg.In, msg.From, msg.Text)
		return chat.Handled
	})

	for {
	}
}
