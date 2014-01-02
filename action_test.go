package cork

import (
	"net/http"
	"testing"
)

var action Action = &DefaultAction{}

func TestUnsupportedAction(t *testing.T) {
	action := action.Create(func() {})
	expect(t, action, nil)
}

func TestHandlerFunc(t *testing.T) {
	hf := func(res http.ResponseWriter, req *http.Request) {}
	action := action.Create(hf)
	_, ok := action.(*handlerAction)
	expect(t, ok, true)
}

type stubHandler struct{}

func (h *stubHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {}

func TestHandler(t *testing.T) {
	action := action.Create(&stubHandler{})
	_, ok := action.(*handlerAction)
	expect(t, ok, true)
}

func TestCorkHandlerFuncAction(t *testing.T) {
	action := action.Create(func(*Response, *Request) {})
	_, ok := action.(Handler)
	expect(t, ok, true)
}
