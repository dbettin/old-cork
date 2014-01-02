package cork

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type stubAction struct{}

func (a *stubAction) Handle(res *Response, req *Request) {}

type stubRouter struct{}

func (r *stubRouter) Route(method string, url *url.URL) *Route {

	if url.String() == "/foo" {
		route := &Route{methodMatch: false, Action: new(stubAction)}
		if method == "GET" {
			route.methodMatch = true
		}
		return route
	}

	return nil
}

func (r *stubRouter) Configure(settings *Cork) {
}

func Test404Dispatch(t *testing.T) {
	d := &defaultDispatcher{Router: new(stubRouter)}
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	d.ServeHTTP(res, req)

	expect(t, res.Code, 404)
}

func Test405Dispatch(t *testing.T) {
	d := &defaultDispatcher{Router: new(stubRouter)}
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/foo", nil)

	d.ServeHTTP(res, req)

	expect(t, res.Code, 405)
}
