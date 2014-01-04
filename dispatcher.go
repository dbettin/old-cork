package cork

import (
	"errors"
	"fmt"
	"net/http"
)

type Dispatcher interface {
	http.Handler
	SettingsHandler
}

type defaultDispatcher struct {
	settings *Cork
	Router
	RequestCreator
	ResponseCreator
}

func (d *defaultDispatcher) Configure(settings *Cork) {
	d.settings = settings
	d.Router = settings.Services.Router
	d.RequestCreator = settings.Services.RequestCreator
	d.ResponseCreator = settings.Services.ResponseCreator
}

func (d *defaultDispatcher) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	response, request := d.setupRequest(res, req)
	defer d.handlePanic(response, request)
	d.dispatch(response, request)
}

func (d *defaultDispatcher) dispatch(response *Response, request *Request) {
	route := d.Route(request.Method, request.URL)

	if route != nil && route.Action != nil {
		d.setup(route, request)
		if route.methodMatch {
			d.runHandlers(route, request)
		} else {
			// call 405 handler
			response.Status(http.StatusMethodNotAllowed)
		}
	} else {
		response.Status(http.StatusNotFound)
	}

}

func (d *defaultDispatcher) handlePanic(response *Response, request *Request) {
	if r := recover(); r != nil {
		// log error
		var err error
		if s, ok := r.(string); ok {
			fmt.Println(s)
			err = errors.New(s)
		} else if e, ok := r.(error); ok {
			err = e
		} else {
			err = errors.New("Unknown error")
		}

		if eh := d.settings.Error; eh != nil {
			request.Error = NewError(err, http.StatusInternalServerError)
			eh.Handle(response, request)
		} else {
			response.Status(http.StatusInternalServerError)
		}
	}
}

func (d *defaultDispatcher) runHandlers(route *Route, request *Request) {
	request.Next()
}

func (d *defaultDispatcher) setupRequest(res http.ResponseWriter, req *http.Request) (*Response, *Request) {
	response := d.NewResponse(res)
	request := d.NewRequest(req, response)
	return response, request
}

func (d *defaultDispatcher) setup(route *Route, request *Request) {
	request.SetRoute(route)
	route.addHandler(route.Action)
}
