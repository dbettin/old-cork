package cork

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	http.ResponseWriter
	settings *Cork
}

type responseWriter struct {
	http.ResponseWriter
}

type ResponseCreator interface {
	NewResponse(http.ResponseWriter, *Cork) *Response
}

type defaultResponseCreator struct{}

func (rc *defaultResponseCreator) NewResponse(res http.ResponseWriter, settings *Cork) *Response {
	response := &Response{}
	response.settings = settings
	response.ResponseWriter = &responseWriter{ResponseWriter: res}
	return response
}

func (r *Response) Status(status int) {
	r.WriteHeader(status)
}

func (r *Response) OK(payload interface{}) {
	r.Send(http.StatusOK, payload)
}

func (r *Response) Error(status int, err error) {
	if r.settings.Error != nil {
	}

}

func (r *Response) Send(status int, payload interface{}) {
	result, err := json.Marshal(payload)
	if err != nil {
		r.Error(http.StatusInternalServerError, err)
	} else {
		r.Status(status)
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
