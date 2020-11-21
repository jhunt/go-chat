package chat

type Handle string
type Context string

type Message struct {
	from Handle
	to   Handle
	in   Context
	text string

	bot Bot
}

func (m Message) Reply(msg string, args ...interface{}) {
	m.bot.Post([]string{string(m.in)}, msg, args...)
}
