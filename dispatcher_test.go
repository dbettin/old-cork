package cork

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type stubAction struct{}

func (a *stubAction) Handle(res *Response, req *Request) {
}

type stubRouter struct{}

func (r *stubRouter) Route(method string, url *url.URL) *Route {

	if url.String() == "/foo" {
		route := &Route{methodMatch: false, Action: new(stubAction)}
		if method == "GET" {
			route.methodMatch = true
		}
		return route
	} else if url.String() == "/panic" {
		panic("error")
	}

	return nil
}

func (r *stubRouter) Configure(settings *Cork) {}

type testError struct {
	Message string
}

var dt *defaultDispatcher

func init() {
	dt = &defaultDispatcher{}
	settings := Pop()
	settings.Services.Router = new(stubRouter)
	dt.Configure(settings)
}

func TestPanicDispatchWithErrorHandler(t *testing.T) {
	dt.settings.Error = HandlerFunc(func(res *Response, req *Request) {
		res.Status(req.Error.Status)
		res.Send(&testError{req.Error.Error()})
	})

	res := sendRequest("GET", "/panic")
	expect(t, res.Code, 500)
	expect(t, res.Body.String(), `{"Message":"error"}`)
}

func TestPanicDispatch(t *testing.T) {
	dt.settings.Error = nil
	res := sendRequest("GET", "/panic")
	expect(t, res.Code, 500)
	expect(t, res.Body.String(), "")
}

func Test404Dispatch(t *testing.T) {
	res := sendRequest("GET", "/foobar")
	expect(t, res.Code, 404)
}

func Test405Dispatch(t *testing.T) {
	res := sendRequest("POST", "/foo")
	expect(t, res.Code, 405)
}

func sendRequest(method string, url string) *httptest.ResponseRecorder {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, nil)
	dt.ServeHTTP(res, req)
	return res
}
