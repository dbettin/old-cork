package cork

import (
	"log"
	"net/url"
	"strings"
)

const (
	GET    string = "GET"
	POST   string = "POST"
	PATCH  string = "PATCH"
	DELETE string = "DELETE"
	PUT    string = "PUT"
)

type Router interface {
	Route(method string, url *url.URL) *Route
	SettingsHandler
}

type defaultRouter struct {
	store routeStore
}

func NewDefaultRouter() *defaultRouter {
	return new(defaultRouter)
}

func (router *defaultRouter) Configure(settings *Cork) {
	router.store = newTrieStore()
	settings.Routes(newDefaultRoutes(settings, router.store))
}

func (router *defaultRouter) Route(method string, url *url.URL) *Route {
	return router.store.find(method, url.String())
}

// TRIE Store

type routeStore interface {
	insert(route *Route)
	find(method string, route string) *Route
}

type trieStore struct {
	root *node
}

func newTrieStore() routeStore {
	return &trieStore{root: &node{segment: &Segment{Name: "/"}}}
}

func (store *trieStore) insert(route *Route) {
	parent := store.root
	if route.isRoot {
		parent.route = route
	} else {
		for _, segment := range route.Segments {
			// find segment in node children
			child := parent.findChild(segment.Name)

			if child == nil {
				child = parent.addChild(segment)
			}
			parent = child
		}
		parent.route = route
	}
}

func (store *trieStore) find(method string, route string) *Route {
	log.Printf("Method: %s, Route: %s", method, route)
	parent := store.root
	//fmt.Printf("Root children length: %v", len(node.children))
	if route != ROOT_ROUTE {
		segments := strings.Split(route, "/")[1:]

		for _, segment := range segments {
			child := parent.findChild(segment)

			if child == nil {
				return nil
			}

			if child.isMatch(method) {
				return child.route
			}
			parent = child
		}
	} else {
		if parent.isMatch(method) {
			return parent.route
		}
	}

	return nil
}

type node struct {
	segment  *Segment
	children []*node
	route    *Route
}

func (parent *node) findChild(pathSegment string) *node {
	if parent.children != nil {
		for _, childNode := range parent.children {

			segment := childNode.segment

			if segment.Variable {
				childNode.segment.Value = pathSegment
				return childNode
			} else if segment.Name == strings.ToLower(pathSegment) {
				return childNode
			}
		}
	}
	return nil
}

func (parent *node) addChild(segment *Segment) (newNode *node) {
	newNode = &node{segment: segment}
	parent.children = append(parent.children, newNode)
	return newNode
}

func (node *node) isMatch(method string) bool {
	if route := node.route; route != nil {
		if route.Method == strings.ToUpper(method) {
			node.route.methodMatch = true
		}
		return true
	}
	return false
}
