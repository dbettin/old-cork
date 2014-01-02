package cork

import (
	"net/http"
	"testing"
)

var store = &stubStore{}
var routes = newDefaultRoutes(Pop(), store)

type stubStore struct {
	route []*Route
}

func (store *stubStore) insert(route *Route) {
	store.route = append(store.route, route)
}

func (store *stubStore) find(method string, route string) *Route {
	return nil
}

func d(response http.ResponseWriter, request *http.Request) {
}

func TestIndexRouteConfig(t *testing.T) {
	routes.Get("/", d)

	if store.route == nil {
		t.Error("Failed to configure route '/'. It was not added to the store.")
	}

	if _, ok := store.route[0].Action.(http.Handler); ok && store.route[0].Template != "/" && store.route[0].Method != GET {
		t.Error("Failed to configure route '/'.")
	}
}

func TestShouldNotAddRouteIfError(t *testing.T) {
	store.route = nil

	routes.Get("thiswillerror/foo", d)

	if store.route != nil {
		t.Error("Improper route should not have been added to store.")
	}
}

func TestShouldNotAddRouteIfWrongAction(t *testing.T) {
	store.route = nil

	routes.Get("/foo", func() {})

	if store.route != nil {
		t.Error("Improper route should not have been added to store.")
	}
}

func TestRouteForwarding(t *testing.T) {
	store.route = nil

	fr := func(r Routes) {
		r.Get("/foo", d)
	}

	routes.Get("/foobar", d)
	routes.Forward("/prefix", fr)
	routes.Get("/boo", d)
	routes.Forward("/prefix2", nil)
	routes.Forward("/prefix2", fr)

	expect(t, store.route[0].Template, "/foobar")
	expect(t, store.route[1].Template, "/prefix/foo")
	expect(t, store.route[2].Template, "/boo")
	expect(t, store.route[3].Template, "/prefix2/foo")
}

type ctx struct{}

func TestAddContext(t *testing.T) {
	store.route = nil

	routes.Get("/foo", d).
		WithContext(&ctx{})

	refute(t, store.route[0].context, nil)
}

func TestAddHandler(t *testing.T) {
	store.route = nil

	hf := func(res *Response, req *Request) {}

	routes.Get("/foo", d).
		HandleFunc(hf)

	routes.Get("/bar", d).
		Handle(HandlerFunc(hf))

	refute(t, store.route[0].Handlers[0], nil)
	refute(t, store.route[1].Handlers[0], nil)
}
