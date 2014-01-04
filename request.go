package cork

import (
	"fmt"
	"net/http"
)

type Request struct {
	*http.Request
	Route   *Route
	Context interface{}
	Params  map[string]string
	Error   *Error

	response    *Response
	nextHandler int
}

func (r *Request) SetRoute(route *Route) {
	r.Route = route
	r.setupParams()
}

func (r *Request) Next() {
	handlers := r.Route.Handlers
	fmt.Println("%v", len(handlers))
	if r.nextHandler < len(handlers) {
		idx := r.nextHandler
		r.nextHandler++
		handlers[idx].Handle(r.response, r)
	}
}

func (r *Request) setupParams() {
	for _, segment := range r.Route.Segments {
		if segment.Variable {
			r.Params[segment.Name] = segment.Value
		}
	}
}

type RequestCreator interface {
	NewRequest(*http.Request, *Response) *Request
}

type defaultRequestCreator struct{}

func (rc *defaultRequestCreator) NewRequest(req *http.Request, res *Response) *Request {
	request := &Request{Request: req, response: res, nextHandler: 0}
	request.Params = make(map[string]string)
	return request
}
