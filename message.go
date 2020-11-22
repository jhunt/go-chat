package chat

type Handle string
type Context string

type Message struct {
	From Handle
	To   Handle
	In   Context
	Text string

	bot Bot
}

func (m Message) Reply(msg string, args ...interface{}) {
	m.bot.Post([]string{string(m.In)}, msg, args...)
}
