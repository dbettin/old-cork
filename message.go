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

func (m *Message) Status(status int) {
	m.Response.WriteHeader(status)
}

func (m *Message) ReplyOK(payload interface{}) {
	m.Reply(http.StatusOK, payload)
}

func (m *Message) ReplyError(status int, err error) {
	if m.Settings.Error != nil {
	}

}

func (m *Message) Reply(status int, payload interface{}) {
	result, err := json.Marshal(payload)
	if err != nil {
		m.ReplyError(http.StatusInternalServerError, err)
	} else {
		m.Status(status)
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
