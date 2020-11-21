package main

import (
	"time"

	"github.com/jhunt/go-chat"
)

func main() {
	bot, _ := chat.Terminal()
	say := func() {
		bot.Post([]string{}, "every 15s: hello, world!")
	}

	bot.On("info", func(msg chat.Message) {
		msg.Reply("no info available at this time")
	})

	t := time.NewTicker(15 * time.Second)
	say()
	for range t.C {
		say()
	}
}
