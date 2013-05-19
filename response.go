package restful

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strings"
)

// If Accept header matching fails, fall back to this type, otherwise
// a "406: Not Acceptable" response is returned.
// Valid values are restful.MIME_JSON and restful.MIME_XML
// Example:
// 	restful.DefaultResponseMimeType = restful.MIME_JSON
var DefaultResponseMimeType string

// Response is a wrapper on the actual http ResponseWriter
// It provides several convenience methods to prepare and write response content.
type Response struct {
	http.ResponseWriter
	accept   string   // content-types what the Http Request says it want to receive
	produces []string // content-types what the Route says it can produce
}

// InternalServerError is a shortcut for .WriteHeader(http.StatusInternalServerError)
func (self Response) InternalServerError() Response {
	self.WriteHeader(http.StatusInternalServerError)
	return self
}

// AddHeader is a shortcut for .Header().Add(header,value)
func (self Response) AddHeader(header string, value string) Response {
	self.Header().Add(header, value)
	return self
}

// WriteEntity marshals the value using the representation denoted by the Accept Header (XML or JSON)
// If no Accept header is specified (or */*) then return the Content-Type as specified by the first in the Route.Produces.
// If an Accept header is specified then return the Content-Type as specified by the first in the Route.Produces that is matched with the Accept header.
// Current implementation ignores any q-parameters in the Accept Header.
func (self Response) WriteEntity(value interface{}) Response {
	if "" == self.accept || "*/*" == self.accept {
		for _, each := range self.produces {
			if MIME_JSON == each {
				self.WriteAsJson(value)
				return self
			}
			if MIME_XML == each {
				self.WriteAsXml(value)
				return self
			}
		}
	} else { // Accept header specified ; scan for each element in Route.Produces
		for _, each := range self.produces {
			if strings.Index(self.accept, each) != -1 {
				if MIME_JSON == each {
					self.WriteAsJson(value)
					return self
				}
				if MIME_XML == each {
					self.WriteAsXml(value)
					return self
				}
			}
		}
	}
	if DefaultResponseMimeType == MIME_JSON {
		self.WriteAsJson(value)
	} else if DefaultResponseMimeType == MIME_XML {
		self.WriteAsXml(value)
	} else {
		self.WriteHeader(http.StatusNotAcceptable)
		self.Write([]byte("406: Not Acceptable"))
	}
	return self
}

// WriteAsXml is a convenience method for writing a value in xml (requires Xml tags on the value)
func (self Response) WriteAsXml(value interface{}) Response {
	output, err := xml.MarshalIndent(value, " ", " ")
	if err != nil {
		self.WriteError(http.StatusInternalServerError, err)
	} else {
		self.Header().Set(HEADER_ContentType, MIME_XML)
		self.Write([]byte(xml.Header))
		self.Write(output)
	}
	return self
}

// WriteAsJson is a convenience method for writing a value in json
func (self Response) WriteAsJson(value interface{}) Response {
	output, err := json.MarshalIndent(value, " ", " ")
	if err != nil {
		self.WriteError(http.StatusInternalServerError, err)
	} else {
		self.Header().Set(HEADER_ContentType, MIME_JSON)
		self.Write(output)
	}
	return self
}

// WriteError is a convenience method for an error status with the actual error
func (self Response) WriteError(status int, err error) Response {
	self.WriteHeader(status)
	if err != nil {
		self.WriteEntity(err.Error())
	}
	return self
}
