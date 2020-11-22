package chat

import (
	"bufio"
	"fmt"
	"os"
)

type TerminalBot struct {
	init  bool
	on    map[string]Handler
	every Handler
}

func Terminal() (Bot, error) {
	return &TerminalBot{
		on: make(map[string]Handler),
	}, nil
}

func (b *TerminalBot) Post(to []string, msg string, args ...interface{}) {
	fmt.Printf(">> "+msg+"\n", args...)
}

func (b *TerminalBot) Every(fn Handler) {
	b.every = fn
	b.listen()
}

func (b *TerminalBot) On(in string, fn Handler) {
	b.on[in] = fn
	b.listen()
}

func (b *TerminalBot) listen() {
	if !b.init {
		b.init = true
		go b.read()
	}
}

func (b *TerminalBot) read() {
	buf := bufio.NewScanner(os.Stdin)
	m := Message{
		From: "",
		To:   "",
		In:   "#stdin",
		Text: buf.Text(),
		bot:  b,
	}

Processing:
	for buf.Scan() {
		if b.every != nil {
			if b.every(m) == Handled {
				continue Processing
			}
		}

		for want, handler := range b.on {
			if want == buf.Text() {
				if handler(m) == Handled {
					continue Processing
				}
			}
		}
	}
}
