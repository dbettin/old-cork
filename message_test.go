package cork

import (
	"net/http"
	"testing"
)

var factory MessageCreator = &defaultMessageCreator{}

func TestNewMessage(t *testing.T) {
	httpreq, _ := http.NewRequest("GET", "/foo", nil)
	req := factory.NewMessage(httpreq, nil, nil)
	refute(t, req, nil)
	refute(t, req.Params, nil)
	refute(t, req.Request, nil)
}

func TestSetRoute(t *testing.T) {
	httpreq, _ := http.NewRequest("GET", "/foo", nil)
	m := factory.NewMessage(httpreq, nil, nil)
	m.SetRoute(&Route{context: "some context"})
	refute(t, m.Route, nil)
	refute(t, m.Context, nil)
}

func TestMessageNext(t *testing.T) {
	httpreq, _ := http.NewRequest("GET", "/foo", nil)
	route := &Route{}
	var counter int
	hf := func(message *Message) {
		counter++
	}

	route.addHandler(HandlerFunc(hf))
	route.addHandler(HandlerFunc(hf))
	req := factory.NewMessage(httpreq, nil, nil)
	req.SetRoute(route)
	expect(t, counter, 0)
	req.Next()
	expect(t, counter, 1)
	req.Next()
	expect(t, counter, 2)
}

func TestMessageParams(t *testing.T) {
	httpreq, _ := http.NewRequest("GET", "/foo", nil)
	route := &Route{}
	route.addSegment(&Segment{Variable: true, Value: "foo", Name: "id"})
	route.addSegment(&Segment{Name: "cork"})
	route.addSegment(&Segment{Variable: true, Value: "bar", Name: "temp"})
	req := factory.NewMessage(httpreq, nil, nil)
	req.SetRoute(route)
	expect(t, len(req.Params), 2)
	expect(t, req.Params["id"], "foo")
	expect(t, req.Params["temp"], "bar")
}
