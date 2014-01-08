package cork

import (
	"encoding/json"
	"net/http"
)

type Message struct {
	*http.Request

	Response http.ResponseWriter
	Route    *Route
	Context  interface{}
	Params   map[string]string
	Error    *Error
	Settings *Cork

	nextHandler int
}

func (m *Message) SetRoute(route *Route) {
	m.Route = route
	m.Context = route.context
	m.setupParams()
}

func (m *Message) Next() {
	handlers := m.Route.Handlers
	if m.nextHandler < len(handlers) {
		idx := m.nextHandler
		m.nextHandler++
		handlers[idx].Handle(m)
	}
}

func (m *Message) setupParams() {
	for _, segment := range m.Route.Segments {
		if segment.Variable {
			m.Params[segment.Name] = segment.Value
		}
	}
}

type MessageCreator interface {
	NewMessage(*http.Request, http.ResponseWriter, *Cork) *Message
}

type defaultMessageCreator struct{}

func (mc *defaultMessageCreator) NewMessage(req *http.Request, res http.ResponseWriter, settings *Cork) *Message {
	message := &Message{Request: req, nextHandler: 0, Settings: settings}
	message.Response = &responseWriter{ResponseWriter: res}
	message.Params = make(map[string]string)
	return message
}

type responseWriter struct {
	http.ResponseWriter
}

func (m *Message) ReturnStatus(status int) {
	m.Response.WriteHeader(status)
}

func (m *Message) ReturnOK(payload interface{}) {
	m.Return(http.StatusOK, payload)
}

func (m *Message) ReturnServerError(err error) {
	m.ReturnError(http.StatusInternalServerError, err)
}

func (m *Message) ReturnError(status int, err error) {
	if me, ok := err.(*Error); ok {
		m.Error = me
	} else {
		m.Error = NewError(err, status)
	}
	m.Settings.Error.Handle(m)
}

func (m *Message) Return(status int, payload interface{}) {
	// this will move to a standard cork handler
	result, err := json.Marshal(payload)
	if err != nil {
		m.ReturnServerError(err)
	} else {
		m.Response.WriteHeader(status)
		m.Response.Header().Set("Content-Type", "application/json")
		m.Response.Write(result)
	}
}

/*
func (rw *responseWriter) WriteHeader(s int) {
	rw.ResponseWriter.WriteHeader(s)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	return size, err
}
*/
