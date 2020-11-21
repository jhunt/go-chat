package chat

import (
	"bufio"
	"fmt"
	"os"
)

type TerminalBot struct {
	init bool
	on   map[string]Handler
}

func Terminal() (Bot, error) {
	return &TerminalBot{
		on: make(map[string]Handler),
	}, nil
}

func (b *TerminalBot) Post(to []string, msg string, args ...interface{}) {
	fmt.Printf(">> "+msg+"\n", args...)
}

func (b *TerminalBot) On(in string, fn Handler) {
	b.on[in] = fn
	if !b.init {
		b.init = true
		go b.read()
	}
}

func (b *TerminalBot) read() {
	buf := bufio.NewScanner(os.Stdin)
	for buf.Scan() {
		for want, handler := range b.on {
			if want == buf.Text() {
				handler(Message{
					from: "",
					to:   "",
					in:   "#stdin",
					text: buf.Text(),
					bot:  b,
				})
				break
			}
		}
	}
}
