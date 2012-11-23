package restful

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strings"
)

const (
	MIME_XML  = "application/xml"
	MIME_JSON = "application/json"
)

// Response is a wrapper on the actual http ResponseWriter
// It provides several convenience methods to prepare and write response content.
type Response struct {
	http.ResponseWriter
	accept string
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

// Marshal the value using the representation denoted by the Accept Header (XML or JSON)
// If no Accept header is specified then return MIME_XML content
func (self Response) Entity(value interface{}) Response {
	if strings.Index(self.accept, MIME_JSON) != -1 {
		self.WriteAsJson(value)
		return self
	}
	//	if strings.Index(self.accept,MIME_XML) != -1 {
	self.WriteAsXml(value)
	return self
}

// Convenience method for writing a value in xml (requires Xml tags on the value)
func (self Response) WriteAsXml(value interface{}) Response {
	output, err := xml.MarshalIndent(value, " ", " ")
	if err != nil {
		self.InternalServerError()
	} else {
		self.Header().Set("Content-Type", MIME_XML)
		self.Write([]byte(xml.Header))
		self.Write(output)
	}
	return self
}

// Convenience method for writing a value in json
func (self Response) WriteAsJson(value interface{}) Response {
	output, err := json.MarshalIndent(value, " ", " ")
	if err != nil {
		self.InternalServerError()
	} else {
		self.Header().Set("Content-Type", MIME_JSON)
		self.Write(output)
	}
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
