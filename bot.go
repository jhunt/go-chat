package chat

type Handler func(Message)

type Bot interface {
	Post([]string, string, ...interface{})
	On(string, Handler)
}
