package cork

import (
	"fmt"
)

type Routes interface {
	Get(template string, action interface{}) HandlerWithExpression
	Forward(prefix string, routes func(r Routes))
}

type HandlerWithExpression interface {
	HandlerExpression
	WithExpression
}

type HandlerExpression interface {
	Handle(Handler) HandlerWithExpression
	HandleFunc(func(*Message)) HandlerWithExpression
}

type WithExpression interface {
	WithContext(interface{})
}

func (routes *defaultRoutes) Get(template string, action interface{}) HandlerWithExpression {
	if routes.prefix != "" {
		template = routes.prefix + template
	}

	route, err := NewRoute(template, GET)
	routes.currentRoute = route
	services := routes.settings.Services

	if err != nil {
		// log here instead
		fmt.Printf("Failed to add route '%s' due to the following error: '%s'", template, err.Error())
	} else {
		action := services.Action.Create(action)
		if action == nil {
			fmt.Printf("Failed to add route '%s' due to the following error: '%s'", template, "Invalid Action.")
		} else {
			routes.currentRoute.Action = action
			routes.store.insert(routes.currentRoute)
		}
	}

	return routes
}

func (routes *defaultRoutes) Forward(prefix string, forwardedRoutes func(r Routes)) {
	if forwardedRoutes != nil {
		routes.prefix = prefix
		defer func() {
			routes.prefix = ""
		}()
		forwardedRoutes(routes)
	}
}

func (routes *defaultRoutes) Handle(handler Handler) HandlerWithExpression {
	routes.currentRoute.addHandler(handler)
	return routes
}

func (routes *defaultRoutes) HandleFunc(handler func(*Message)) HandlerWithExpression {
	routes.currentRoute.addHandler(HandlerFunc(handler))
	return routes
}

func (routes *defaultRoutes) WithContext(context interface{}) {
	routes.currentRoute.addContext(context)
}

type defaultRoutes struct {
	currentRoute *Route
	prefix       string
	store        routeStore
	settings     *Cork
}

func newDefaultRoutes(settings *Cork, store routeStore) Routes {
	return &defaultRoutes{settings: settings, store: store}
}
