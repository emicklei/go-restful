package restful

import (
	"net/http"
)

// Response is a wrapper on the actual http ResponseWriter
// It provides several convenience methods to prepare responses.
type Response struct {
	http.ResponseWriter
}

// Shortcut for .WriteHeader(http.StatusInternalServerError)
func (self Response) InternalServerError() Response {
	self.WriteHeader(http.StatusInternalServerError)
	return self
}

// Shortcut for .Header().Add(header,value)
func (self Response) AddHeader(header string, value string) Response {
	self.Header().Add(header, value)
	return self
}
func (self Response) Entity(entity interface{}) Response {
	return self
}

// Flush and close the underlying ResponseWriter
// From https://github.com/nharbour/web.go/blob/master/web.go
//func (self Response) Close() {
//    rwc, buf, _ := self.(http.Hijacker).Hijack()
//    if buf != nil {
//        buf.Flush()
//    }
//
//    if rwc != nil {
//        rwc.Close()
//    }
//}
