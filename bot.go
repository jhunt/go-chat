package chat

type Then int

const (
	Handled Then = iota
	Continue
)

type Handler func(Message, ...string) Then

type Bot interface {
	Post([]string, string, ...interface{})
	Every(Handler)
	On(string, Handler)
}
