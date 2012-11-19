package restful

import (
	"net/http"
)

type Response struct {
	//Writer http.ResponseWriter
	responseWriter
}

func (self Response) StatusOK() Response {
	self.WriteHeader(http.StatusOK)
	return self
}
func (self Response) StatusError() Response {
	self.WriteHeader(http.StatusInternalServerError)
	return self
}
func (self Response) AddHeader(header string, value string) Response {
	return self
}
func (self Response) Entity(entity interface{}) Response {
	return self
}

// From https://github.com/nharbour/web.go/blob/master/web.go
type responseWriter interface {
    Header() http.Header
    WriteHeader(status int)
    Write(data []byte) (n int, err error)
    Close()
}
type responseWriter struct {
    http.ResponseWriter
}
func (c *responseWriter) Close() {
    rwc, buf, _ := c.ResponseWriter.(http.Hijacker).Hijack()
    if buf != nil {
        buf.Flush()
    }

    if rwc != nil {
        rwc.Close()
    }
}