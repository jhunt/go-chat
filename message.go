package chat

type Handle string
type Context string

type Message struct {
	from Handle
	to   Handle
	in   Context

	bot Bot
}

func (m Message) Reply(msg string, args ...interface{}) {
	m.bot.Post([]string{}, msg, args...)
}
