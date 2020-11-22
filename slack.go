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
	init  bool
	on    map[string]Handler
	every Handler
	c     slack.Client
}

func Slack(token string) (Bot, error) {
	c, err := slack.Connect(token)
	if err != nil {
		return nil, err
	}

	return &SlackBot{
		c:  c,
		on: make(map[string]Handler),
	}, nil
}

func (b *SlackBot) Post(to []string, msg string, args ...interface{}) {
	text := fmt.Sprintf(msg, args...)

	for _, whom := range to {
		fmt.Printf("saying '%s' to %s\n", text, whom)
		b.c.Send(slack.Message{
			Type:    "message",
			Channel: whom,
			Text:    text,
		})
	}
}

func (b *SlackBot) Every(fn Handler) {
	b.every = fn
	b.listen()
}

func (b *SlackBot) On(in string, fn Handler) {
	b.on[in] = fn
	b.listen()
}

func (b *SlackBot) listen() {
	if !b.init {
		b.init = true
		go b.read()
	}
}

func (b *SlackBot) read() {
Processing:
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
			From:     Handle(m.User),
			In:       Context(m.Channel),
			Text:     m.Text,
			bot:      b,
		}

		if b.every != nil {
			if b.every(msg) == Handled {
				continue Processing
			}
		}

		if m.IsDirected(b.c.Name) {
			msg.Text = salutation.ReplaceAllString(m.Text, "")

			for want, handler := range b.on {
				if want == m.Text {
					if handler(msg) == Handled {
						continue Processing
					}
				}
			}
		}
	}
}
