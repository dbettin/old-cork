package cork

import (
	"net/http"
	"testing"
)

func TestNewRequest(t *testing.T) {
	httpreq, _ := http.NewRequest("GET", "/foo", nil)
	req := NewRequest(&Route{}, httpreq, nil)
	refute(t, req, nil)
	refute(t, req.Route, nil)
	refute(t, req.Request, nil)
}

func TestRequestNext(t *testing.T) {
	httpreq, _ := http.NewRequest("GET", "/foo", nil)
	route := &Route{}
	var counter int
	hf := func(res *Response, req *Request) {
		counter++
	}

	route.addHandler(HandlerFunc(hf))
	route.addHandler(HandlerFunc(hf))
	req := NewRequest(route, httpreq, nil)
	expect(t, counter, 0)
	req.Next()
	expect(t, counter, 1)
	req.Next()
	expect(t, counter, 2)
}

func TestRequestParams(t *testing.T) {
	httpreq, _ := http.NewRequest("GET", "/foo", nil)
	route := &Route{}
	route.addSegment(&Segment{Variable: true, Value: "foo", Name: "id"})
	route.addSegment(&Segment{Name: "cork"})
	route.addSegment(&Segment{Variable: true, Value: "bar", Name: "temp"})
	req := NewRequest(route, httpreq, nil)
	expect(t, len(req.Params), 2)
	expect(t, req.Params["id"], "foo")
	expect(t, req.Params["temp"], "bar")
}
