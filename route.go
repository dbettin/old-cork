package cork

import (
	"errors"
	"regexp"
	"strings"
)

const ROOT_ROUTE = "/"

var compiled_var_regex, _ = regexp.Compile(`^{(\w*)}$`)

type Route struct {
	Template string
	Method   string
	Action   Handler
	Segments []*Segment
	Handlers []Handler

	methodMatch bool
	context     interface{}
	isRoot      bool
}

type Segment struct {
	Name     string
	Variable bool
	Value    string
}

func NewRoute(template string, method string) (*Route, error) {
	route := &Route{Template: template, Method: method}
	route.isRoot = checkIfRoot(template)

	if route.isRoot {
		route.addSegment(&Segment{Name: template})
	} else {
		segments := strings.Split(template, "/")

		if segments[0] != "" {
			return nil, errors.New("Please specify a full path for route. Missing a leading '/' for: " + template)
		}

		segments = segments[1:]
		for _, name := range segments {
			route.parseSegment(name)
		}
	}

	return route, nil
}

func (route *Route) addHandler(handler Handler) {
	if handler != nil {
		route.Handlers = append(route.Handlers, handler)
	}
}

func (route *Route) addContext(context interface{}) {
	if context != nil {
		route.context = context
	}
}

func (route *Route) parseSegment(name string) {
	match := compiled_var_regex.FindStringSubmatch(name)
	segment := &Segment{Name: strings.ToLower(name)}

	if len(match) > 0 && match[0] != "" {
		segment.Name = strings.ToLower(match[1])
		segment.Variable = true
	}
	route.addSegment(segment)
}

func (route *Route) addSegment(segment *Segment) {
	route.Segments = append(route.Segments, segment)
}

func checkIfRoot(path string) bool {
	if path == ROOT_ROUTE {
		return true
	}
	return false
}
