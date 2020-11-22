package chat

import (
	"fmt"
	"os"
	"regexp"

	"github.com/jhunt/go-chat/lib/slack"
)

var salutation *regexp.Regexp

func init() {
	salutation = regexp.MustCompile(`^\s*<@.*?>\s*,?\s*`)
}

type SlackBot struct {
	init bool
	on map[string]Handler
	c slack.Client
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

func (b *SlackBot) On(in string, fn Handler) {
	b.on[in] = fn
	if !b.init {
		b.init = true
		go b.read()
	}
}

func (b *SlackBot) read() {
	for {
		m, err := b.c.Receive()
		if err != nil {
			fmt.Fprintf(os.Stderr, "oops: %s\n", err)
			continue
		}

		fmt.Fprintf(os.Stderr, "recv: %s\n", m)
		if m.Type != "message" || !m.IsDirected(b.c.Name) {
			continue
		}
		m.Text = salutation.ReplaceAllString(m.Text, "")

		fmt.Fprintf(os.Stderr, "[%s]\n", m.Text)
		for want, handler := range b.on {
			if want == m.Text {
				fmt.Fprintf(os.Stderr, "invoking handler...\n")
				handler(Message{
					from: "", // FIXME
					to:   "", // FIXME
					in:   Context(m.Channel),
					text: m.Text,
					bot:  b,
				})
			}
		}
	}
}
