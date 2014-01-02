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

	response    *Response
	nextHandler int
}

type RequestCreator interface {
	NewRequest(*Route, *http.Request, *Response) *Request
}

type defaultRequestCreator struct{}

func (rc *defaultRequestCreator) NewRequest(route *Route, req *http.Request, res *Response) *Request {
	request := &Request{Request: req, Route: route,
		Context: route.context, response: res, nextHandler: 0}
	request.Params = make(map[string]string)
	request.setupParams()
	return request
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
