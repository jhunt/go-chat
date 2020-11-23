package main

import (
	"fmt"
	"os"
	"regexp"

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

	lower := regexp.MustCompile(`[a-z]`)
	bot.On(`\s*is\s+(.*)\s+upper\s*case\?\s*$`,
		func(msg chat.Message, args ...string) chat.Then {
			if !lower.MatchString(args[1]) {
				msg.Reply("yup, looks like '%s' upper case alright", args[1])
			} else {
				msg.Reply("nope; '%s' is definitely not uppercase", args[1])
			}
			return chat.Handled
		})

	for {
	}
}
