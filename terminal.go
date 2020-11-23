package chat

import (
	"regexp"
	"bufio"
	"fmt"
	"os"
)

type TerminalBot struct {
	init  bool
	on    map[*regexp.Regexp]Handler
	every Handler
}

func Terminal() (Bot, error) {
	return &TerminalBot{
		on: make(map[*regexp.Regexp]Handler),
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
	b.on[regexp.MustCompile(in)] = fn
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
Processing:
	for buf.Scan() {
		msg := Message{
			From: "",
			To:   "",
			In:   "#stdin",
			Text: buf.Text(),
			bot:  b,
		}

		if b.every != nil {
			if b.every(msg) == Handled {
				continue Processing
			}
		}

		for pat, handler := range b.on {
			if matches := pat.FindStringSubmatch(msg.Text); matches != nil {
				if handler(msg, matches...) == Handled {
					continue Processing
				}
			}
		}
	}
}
