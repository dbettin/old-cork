package cork

import (
	"reflect"
	"testing"
)

func expect(t *testing.T, value interface{}, expected interface{}) {
	if value != expected {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", expected, reflect.TypeOf(expected), value, reflect.TypeOf(value))
	}
}

func refute(t *testing.T, value interface{}, notexpected interface{}) {
	if value == notexpected {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", notexpected, reflect.TypeOf(notexpected),
			value, reflect.TypeOf(value))
	}
}
