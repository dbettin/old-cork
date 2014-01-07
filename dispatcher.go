package cork

import (
	"errors"
	"net/http"
)

type Dispatcher interface {
	http.Handler
	SettingsHandler
}

type defaultDispatcher struct {
	settings *Cork
	Router
	MessageCreator
}

func (d *defaultDispatcher) Configure(settings *Cork) {
	d.settings = settings
	d.Router = settings.Services.Router
	d.MessageCreator = settings.Services.MessageCreator
}

func (d *defaultDispatcher) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	message := d.setupMessage(res, req)
	defer d.handlePanic(message)
	d.dispatch(message)
}

func (d *defaultDispatcher) dispatch(message *Message) {
	route := d.Route(message.Method, message.URL)

	if route != nil && route.Action != nil {
		d.setup(route, message)
		if route.methodMatch {
			d.runHandlers(route, message)
		} else {
			// call 405 handler
			message.Status(http.StatusMethodNotAllowed)
		}
	} else {
		message.Status(http.StatusNotFound)
	}

}

func (d *defaultDispatcher) handlePanic(message *Message) {
	if r := recover(); r != nil {
		// log error
		message.Status(http.StatusInternalServerError)

		var err error
		if s, ok := r.(string); ok {
			err = errors.New(s)
		} else if e, ok := r.(error); ok {
			err = e
		} else {
			err = errors.New("Unknown error")
		}

		if eh := d.settings.Error; eh != nil {
			message.Error = NewError(err, http.StatusInternalServerError)
			eh.Handle(message)
		}
	}
}

func (d *defaultDispatcher) runHandlers(route *Route, message *Message) {
	message.Next()
}

func (d *defaultDispatcher) setupMessage(res http.ResponseWriter, req *http.Request) *Message {
	return d.NewMessage(req, res, d.settings)
}

func (d *defaultDispatcher) setup(route *Route, message *Message) {
	message.SetRoute(route)
	route.addHandler(route.Action)
}
