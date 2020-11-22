package main

import (
	"fmt"
	"html"
	"os"
	"time"

	"github.com/jhunt/go-chat"
	"github.com/jhunt/go-rss"
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
		feed, err := rss.Fetch("https://www.motherjones.com/feed/")
		if err != nil {
			bot.Post(channels, "... uh... **%s**\n", err)
			return
		}
		for _, item := range feed.Channel.Items[0:3] {
			bot.Post(channels, "%s\n", item.Link)
			bot.Post(channels, "> %s\n\n\n", html.UnescapeString(item.Description))
		}
	}

	t := time.NewTicker(1 * time.Hour)
	say()
	for range t.C {
		say()
	}
}
