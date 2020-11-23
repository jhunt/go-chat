package chat

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type TerminalBot struct {
	init bool
	d    Dispatcher
}

func Terminal() (Bot, error) {
	return &TerminalBot{}, nil
}

func (b *TerminalBot) Post(to []string, msg string, args ...interface{}) {
	fmt.Printf(">> "+msg+"\n", args...)
}

func (b *TerminalBot) Every(fn Handler) {
	b.d.Every(fn)
	b.listen()
}

func (b *TerminalBot) On(in string, fn Handler) {
	b.d.On(in, fn)
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
	for buf.Scan() {
		msg := Message{
			Received:  time.Now(),
			Addressed: true,

			From: "",
			To:   "",
			In:   "#stdin",
			Text: buf.Text(),
			bot:  b,
		}

		b.d.Dispatch(msg)
	}
}
