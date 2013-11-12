package restful

// Copyright 2013 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

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
	accept        string   // content-types what the Http Request says it want to receive
	produces      []string // content-types what the Route says it can produce
	statusCode    int      // HTTP status code that has been written explicity (if zero then net/http has written 200)
	contentLength int      // number of bytes written for the response body
}

func newResponse(httpWriter http.ResponseWriter) *Response {
	return &Response{httpWriter, "", []string{}, http.StatusOK, 0} // empty content-types
}

// InternalServerError writes the StatusInternalServerError header.
// DEPRECATED, use r.WriteHeader(http.StatusInternalServerError)
func (r Response) InternalServerError() Response {
	r.WriteHeader(http.StatusInternalServerError)
	return r
}

// AddHeader is a shortcut for .Header().Add(header,value)
func (r Response) AddHeader(header string, value string) Response {
	r.Header().Add(header, value)
	return r
}

// WriteEntity marshals the value using the representation denoted by the Accept Header (XML or JSON)
// If no Accept header is specified (or */*) then return the Content-Type as specified by the first in the Route.Produces.
// If an Accept header is specified then return the Content-Type as specified by the first in the Route.Produces that is matched with the Accept header.
// Current implementation ignores any q-parameters in the Accept Header.
func (r *Response) WriteEntity(value interface{}) *Response {
	if "" == r.accept || "*/*" == r.accept {
		for _, each := range r.produces {
			if MIME_JSON == each {
				r.WriteAsJson(value)
				return r
			}
			if MIME_XML == each {
				r.WriteAsXml(value)
				return r
			}
		}
	} else { // Accept header specified ; scan for each element in Route.Produces
		for _, each := range r.produces {
			if strings.Index(r.accept, each) != -1 {
				if MIME_JSON == each {
					r.WriteAsJson(value)
					return r
				}
				if MIME_XML == each {
					r.WriteAsXml(value)
					return r
				}
			}
		}
	}
	if DefaultResponseMimeType == MIME_JSON {
		r.WriteAsJson(value)
	} else if DefaultResponseMimeType == MIME_XML {
		r.WriteAsXml(value)
	} else {
		r.WriteHeader(http.StatusNotAcceptable)
		r.Write([]byte("406: Not Acceptable"))
	}
	return r
}

// WriteAsXml is a convenience method for writing a value in xml (requires Xml tags on the value)
func (r *Response) WriteAsXml(value interface{}) *Response {
	output, err := xml.MarshalIndent(value, " ", " ")
	if err != nil {
		r.WriteError(http.StatusInternalServerError, err)
	} else {
		r.Header().Set(HEADER_ContentType, MIME_XML)
		r.Write([]byte(xml.Header))
		r.Write(output)
	}
	return r
}

// WriteAsJson is a convenience method for writing a value in json
func (r *Response) WriteAsJson(value interface{}) *Response {
	output, err := json.MarshalIndent(value, " ", " ")
	if err != nil {
		r.WriteError(http.StatusInternalServerError, err)
	} else {
		r.Header().Set(HEADER_ContentType, MIME_JSON)
		r.Write(output)
	}
	return r
}

// WriteError write the http status and the error string on the response.
// DEPRECATED; use WriteErrorString(status,reason)
func (r *Response) WriteError(httpStatus int, err error) *Response {
	return r.WriteErrorString(httpStatus, err.Error())
}

// WriteServiceError is a convenience method for a responding with a ServiceError and a status
func (r *Response) WriteServiceError(httpStatus int, err ServiceError) *Response {
	r.WriteHeader(httpStatus)
	r.WriteEntity(err)
	return r
}

// WriteErrorString is a convenience method for an error status with the actual error
func (r *Response) WriteErrorString(status int, errorReason string) *Response {
	r.WriteHeader(status)
	r.Write([]byte(errorReason))
	return r
}

// WriteHeader is overridden to remember the Status Code that has been written.
func (r *Response) WriteHeader(httpStatus int) {
	r.statusCode = httpStatus
	r.ResponseWriter.WriteHeader(httpStatus)
}

// StatusCode returns the code that has been written using WriteHeader.
func (r Response) StatusCode() int {
	if 0 == r.statusCode {
		// no status code has been written yet; assume OK
		return http.StatusOK
	}
	return r.statusCode
}

// Write writes the data to the connection as part of an HTTP reply.
// Write is part of http.ResponseWriter interface
func (r *Response) Write(bytes []byte) (int, error) {
	written, err := r.ResponseWriter.Write(bytes)
	r.contentLength += written
	return written, err
}

// ContentLength returns the number of bytes written for the response content
func (r Response) ContentLength() int {
	return r.contentLength
}
