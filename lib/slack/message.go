package slack

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type Message struct {
	TS       string    `json:"ts,omitempty"`
	Received time.Time `json:"-"`

	ID      uint64 `json:"id"`
	Type    string `json:"type"`
	User    string `json:"user"`
	Channel string `json:"channel"`
	Text    string `json:"text"`

	interned bool
}

func (m Message) IsDirected(at string) bool {
	return strings.HasPrefix(m.Channel, "D") || strings.HasPrefix(m.Text, "<@"+at+">")
}

func (m Message) String() string {
	return fmt.Sprintf("[%d] %s: %s> %s", m.ID, m.Type, m.Channel, m.Text)
}

func prefix(sigil, name string) string {
	if strings.HasPrefix(name, sigil) {
		return name
	}
	return sigil + name
}

func (c Client) intern(m *Message) {
	if m.interned {
		return
	}
	if m.User != "" {
		m.User = c.name2id(prefix("@", m.User))
	}
	if m.Channel != "" && !strings.HasPrefix(m.Channel, "D") {
		m.Channel = c.name2id(prefix("#", m.Channel))
	}

	re := regexp.MustCompile(`<(.*?)>`)
	m.Text = replace(re, m.Text, func(match []string) []string {
		for i, name := range match {
			match[i] = c.name2id(name)
		}
		return match
	})
}

func (c Client) extern(m *Message) {
	if !m.interned {
		return
	}
	if m.User != "" {
		m.User = c.id2name(m.User)
	}
	if m.Channel != "" {
		m.Channel = c.id2name(m.Channel)
	}

	m.Text = regexp.MustCompile(`<[@#](.*?)(\|.*?)?>`).ReplaceAllString(m.Text, "<$1>")

	re := regexp.MustCompile(`<(.*?)>`)
	m.Text = replace(re, m.Text, func(match []string) []string {
		for i, id := range match {
			match[i] = c.id2name(id)
		}
		return match
	})
}
