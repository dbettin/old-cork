package cork

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	http.ResponseWriter
}

type responseWriter struct {
	http.ResponseWriter
}

type ResponseCreator interface {
	NewResponse(http.ResponseWriter) *Response
}

type defaultResponseCreator struct{}

func (rc *defaultResponseCreator) NewResponse(res http.ResponseWriter) *Response {
	response := &Response{}
	response.ResponseWriter = &responseWriter{ResponseWriter: res}
	return response
}

func (r *Response) Send(payload interface{}) {
	result, err := json.Marshal(payload)

	if err != nil {
		// to do - handle errors in global manner
		fmt.Println("Error duing send")

	} else {
		r.Header().Set("Content-Type", "application/json")
		r.Write(result)
	}
}

/*
func (rw *responseWriter) WriteHeader(s int) {
	rw.ResponseWriter.WriteHeader(s)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	return size, err
}
*/
