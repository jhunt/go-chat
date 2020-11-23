package main

import (
	"fmt"
	"os"
	"time"

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

	say := func() {
		bot.Post(channels, "every 15s: hello, world!")
	}

	bot.On("info", func(msg chat.Message, _ ...string) chat.Then {
		msg.Reply("no info available at this time")
		return chat.Handled
	})

	t := time.NewTicker(15 * time.Second)
	say()
	for range t.C {
		say()
	}
}
