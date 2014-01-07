package cork

type Handler interface {
	Handle(*Message)
}

type HandlerFunc func(*Message)

func (h HandlerFunc) Handle(message *Message) {
	h(message)
}
