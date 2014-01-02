package cork

type Handler interface {
	Handle(*Response, *Request)
}

type HandlerFunc func(*Response, *Request)

func (h HandlerFunc) Handle(res *Response, req *Request) {
	h(res, req)
}
