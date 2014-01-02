package cork

import (
	"strings"
	"testing"
)

func TestIndex(t *testing.T) {
	route, _ := NewRoute("/", GET)

	if !route.isRoot {
		t.Error("Index route is not a root path.")
	}

	if len(route.Segments) <= 0 {
		t.Error("Incorrect number of segments for index route")
	}

	if route.Method != "GET" {
		t.Error("Incorrect method was set.")
	}
}

func TestSingleLiteral(t *testing.T) {
	route, _ := NewRoute("/literalpath", GET)

	if route.isRoot {
		t.Error("Should not be a root template.")
	}

	if len(route.Segments) != 1 {
		t.Errorf("Should only have one segment but has %v segments.", len(route.Segments))
	}

	if route.Segments[0].Name != "literalpath" {
		t.Errorf("Segment should not have name %v", route.Segments[0].Name)
	}
}

func TestSegments(t *testing.T) {
	route, _ := NewRoute("/foo/bar/cheese", GET)

	if len(route.Segments) != 3 {
		t.Errorf("Should only have three segments but has %v segments.", len(route.Segments))
	}

	if route.Segments[0].Name != "foo" {
		t.Errorf("Should have 'foo' as first segment but has '%s' instead.", route.Segments[0].Name)
	}
}

func TestInvalidPathSeparator(t *testing.T) {
	_, err := NewRoute("foo/cork/open", GET)

	if err != nil && !strings.HasPrefix(err.Error(), "Please specify a full path for route") {
		t.Errorf("Should have failed due to missing leading slash.")
	}
}

func TestSimpleVariableMatch(t *testing.T) {
	route, _ := NewRoute("/foo/{bar}", GET)

	if route == nil {
		t.Error("Should have parsed simple variable path.")
	}

	segment := route.Segments[1]
	if !segment.Variable || segment.Name != "bar" {
		t.Errorf("Should have variable named 'bar'. Segment is %v", segment.Variable)
	}
}
