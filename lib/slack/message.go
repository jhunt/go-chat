package slack

import (
	"fmt"
	"strings"
)

type Message struct {
	ID      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func (m Message) IsDirected(at string) bool {
	return strings.HasPrefix(m.Channel, "D") || strings.HasPrefix(m.Text, "<@"+at+">")
}

func (m Message) String() string {
	return fmt.Sprintf("[%d] %s: %s> %s", m.ID, m.Type, m.Channel, m.Text)
}
