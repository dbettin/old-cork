package cork

import (
	"testing"
)

func TestDontFindRoot(t *testing.T) {
	store := newTrieStore()
	route := store.find(GET, "/")

	assertNotFound(t, route, "/")
}

func TestFindRoot(t *testing.T) {
	store := newTrieStore()
	entry, _ := NewRoute("/", GET)
	entry.isRoot = true
	store.insert(entry)
	route := store.find(GET, "/")

	assertFound(t, route, "/")
}

func TestFindNonRoot(t *testing.T) {
	route := insertFind(GET, "/foo")

	assertFound(t, route, "/foo")
}

func TestFindNonRootWithMultipleRoutes(t *testing.T) {
	store := insert(GET, "/foo")
	insertWithStore(store, GET, "/foo2")
	insertWithStore(store, GET, "/foo3")
	route := store.find(GET, "/foo2")

	assertFound(t, route, "/foo2")
}

func TestMultipleSegments(t *testing.T) {
	store := insert(GET, "/foo/bar")
	insertWithStore(store, GET, "/foo2/bar2")
	insertWithStore(store, GET, "/foo/bar3")
	route := store.find(GET, "/foo/bar3")

	assertFound(t, route, "/foo/bar3")
}

func TestDontFindNonRoot(t *testing.T) {
	store := newTrieStore()
	route := store.find(GET, "/nofoo")

	assertNotFound(t, route, "/nofoo")
}

func TestWrongMethod(t *testing.T) {
	store := insert(GET, "/foo")
	route := store.find(POST, "/foo")

	assertNotFound(t, route, "/foo")
}

func TestFindVariablePath(t *testing.T) {
	store := insert(GET, "/foo/{bar}")
	route := store.find(GET, "/foo/cork")

	assertFound(t, route, "/foo/cork")

	if route.Segments[1].Value != "cork" {
		t.Errorf("Should have value '%s' for variable 'cork'", route.Segments[1].Value)
	}
}

func TestCaseInsensitiveFind(t *testing.T) {
	route := insertFind(GET, "/FoO")
	assertFound(t, route, "/FoO")
}

// helpers

func assertNotFound(t *testing.T, route *Route, path string) {
	if route != nil && route.methodMatch {
		t.Errorf("Should not have found '%s' route since it wasn't added to the store.", path)
	}
}

func assertFound(t *testing.T, route *Route, path string) {
	if route == nil {
		t.Errorf("Should have found '%s' route since it was added to the store", path)
	}
}

func insertFind(method string, path string) *Route {
	store := insert(method, path)
	route := store.find(GET, path)
	return route
}

func insert(method string, path string) routeStore {
	store := newTrieStore()
	insertWithStore(store, method, path)
	return store
}

func insertWithStore(store routeStore, method string, path string) {
	entry, _ := NewRoute(path, method)
	store.insert(entry)
}
