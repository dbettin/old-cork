package cork

import (
	"net/http"
)

type Action interface {
	Create(action interface{}) Handler
}

type DefaultAction struct{}

func (a DefaultAction) Create(action interface{}) Handler {

	switch action := action.(type) {
	case http.Handler:
		return &handlerAction{handler: action}
	case func(http.ResponseWriter, *http.Request):
		return &handlerAction{handler: http.HandlerFunc(action)}
	case func(*Message):
		return HandlerFunc(action)
	default:
		return nil
	}
}

type handlerAction struct {
	handler http.Handler
}

func (a *handlerAction) Handle(message *Message) {
	a.handler.ServeHTTP(message.Response, message.Request)
}
