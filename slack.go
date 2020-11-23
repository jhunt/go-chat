package chat

import (
	"fmt"
	"regexp"

	"github.com/jhunt/go-chat/lib/slack"
)

var salutation *regexp.Regexp

func init() {
	salutation = regexp.MustCompile(`^\s*<@.*?>\s*,?\s*`)
}

type SlackBot struct {
	init bool
	c    slack.Client
	d    Dispatcher
}

func Slack(token string) (Bot, error) {
	c, err := slack.Connect(token)
	if err != nil {
		return nil, err
	}

	return &SlackBot{c: c}, nil
}

func (b *SlackBot) Post(to []string, msg string, args ...interface{}) {
	text := fmt.Sprintf(msg, args...)

	for _, whom := range to {
		b.c.Send(slack.Message{
			Type:    "message",
			Channel: whom,
			Text:    text,
		})
	}
}

func (b *SlackBot) Every(fn Handler) {
	b.d.Every(fn)
	b.listen()
}

func (b *SlackBot) On(in string, fn Handler) {
	b.d.On(in, fn)
	b.listen()
}

func (b *SlackBot) listen() {
	if !b.init {
		b.init = true
		go b.read()
	}
}

func (b *SlackBot) read() {
	for {
		m, err := b.c.Receive()
		if err != nil {
			continue
		}

		if m.Type != "message" {
			continue
		}

		msg := Message{
			Received: m.Received,

			From: Handle(m.User),
			In:   Context(m.Channel),
			Text: m.Text,
			bot:  b,
		}
		if m.IsDirected(b.c.Name) {
			msg.Addressed = true
			msg.Text = salutation.ReplaceAllString(m.Text, "")
		}

		b.d.Dispatch(msg)
	}
}
