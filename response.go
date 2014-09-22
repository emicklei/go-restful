package restful

// Copyright 2013 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

import (
	"net/http"
	"strings"
)

// DEPRECATED, use DefaultResponseContentType(mime)
var DefaultResponseMimeType string

//PrettyPrintResponses controls the indentation feature of XML and JSON
//serialization in the response methods WriteEntity, WriteAsJson, and
//WriteAsXml.
var PrettyPrintResponses = true

// Response is a wrapper on the actual http ResponseWriter
// It provides several convenience methods to prepare and write response content.
type Response struct {
	http.ResponseWriter
	requestAccept string   // mime-type what the Http Request says it wants to receive
	routeProduces []string // mime-types what the Route says it can produce
	statusCode    int      // HTTP status code that has been written explicity (if zero then net/http has written 200)
	contentLength int      // number of bytes written for the response body
}

// Creates a new response based on a http ResponseWriter.
func NewResponse(httpWriter http.ResponseWriter) *Response {
	return &Response{httpWriter, "", []string{}, http.StatusOK, 0} // empty content-types
}

// If Accept header matching fails, fall back to this type, otherwise
// a "406: Not Acceptable" response is returned.
// Valid values are restful.MIME_JSON and restful.MIME_XML
// Example:
// 	restful.DefaultResponseContentType(restful.MIME_JSON)
func DefaultResponseContentType(mime string) {
	DefaultResponseMimeType = mime
}

// InternalServerError writes the StatusInternalServerError header.
// DEPRECATED, use WriteErrorString(http.StatusInternalServerError,reason)
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
// If the value is nil then nothing is written. You may want to call WriteHeader(http.StatusNotFound) instead.
// Current implementation ignores any q-parameters in the Accept Header.
func (r *Response) WriteEntity(value interface{}) error {
	if value == nil { // do not write a nil representation
		return nil
	}
	var e EntityEncoder
	for _, qualifiedMime := range strings.Split(r.requestAccept, ",") {
		mime := strings.Trim(strings.Split(qualifiedMime, ";")[0], " ")
		if 0 == len(mime) || mime == "*/*" {
			for _, each := range r.routeProduces {
				e = EntityEncoderForMIME(each)
				if e != nil {
					return r.WriteEntityThroughEncoder(value, e)
				}
			}
		} else { // mime is not blank; see if we have a match in Produces
			for _, each := range r.routeProduces {
				if mime == each {
					e = EntityEncoderForMIME(each)
					if e != nil {
						return r.WriteEntityThroughEncoder(value, e)
					}
				}
			}
		}
	}
	if DefaultResponseMimeType == MIME_JSON {
		return r.WriteAsJson(value)
	} else if DefaultResponseMimeType == MIME_XML {
		return r.WriteAsXml(value)
	} else {
		r.WriteHeader(http.StatusNotAcceptable) // for recording only
		r.ResponseWriter.WriteHeader(http.StatusNotAcceptable)
		if _, err := r.Write([]byte("406: Not Acceptable")); err != nil {
			return err
		}
	}
	return nil
}

// WriteEntityThroughEncoder writes a value through the provided EntityEncoder
func (r *Response) WriteEntityThroughEncoder(value interface{}, e EntityEncoder) error {
	var output []byte
	var err error

	if value == nil {
		// do not write a nil representation
		return nil
	}
	e.SetResponse(r)
	if PrettyPrintResponses {
		output, err = e.MarshalIndent(value, " ", " ")
	} else {
		output, err = e.Marshal(value)
	}
	if err != nil {
		return r.WriteError(http.StatusInternalServerError, err)
	}
	r.Header().Set(HEADER_ContentType, e.MIME())
	if r.statusCode > 0 { // a WriteHeader was intercepted
		r.ResponseWriter.WriteHeader(r.statusCode)
	}
	if _, err = r.Write(output); err != nil {
		return err
	}
	return nil
}

// WriteAsXml is a convenience method for writing a value in xml (requires Xml tags on the value)
func (r *Response) WriteAsXml(value interface{}) error {
	e := EntityEncoderForMIME(MIME_XML)
	if e == nil {
		panic("builtin XML encoder not registered")
	}
	return r.WriteEntityThroughEncoder(value, e)
}

// WriteAsJson is a convenience method for writing a value in json
func (r *Response) WriteAsJson(value interface{}) error {
	e := EntityEncoderForMIME(MIME_JSON)
	if e == nil {
		panic("builtin JSON encoder not registered")
	}
	return r.WriteEntityThroughEncoder(value, e)
}

// WriteError write the http status and the error string on the response.
func (r *Response) WriteError(httpStatus int, err error) error {
	return r.WriteErrorString(httpStatus, err.Error())
}

// WriteServiceError is a convenience method for a responding with a ServiceError and a status
func (r *Response) WriteServiceError(httpStatus int, err ServiceError) error {
	r.WriteHeader(httpStatus) // for recording only
	return r.WriteEntity(err)
}

// WriteErrorString is a convenience method for an error status with the actual error
func (r *Response) WriteErrorString(status int, errorReason string) error {
	r.statusCode = status // for recording only
	r.ResponseWriter.WriteHeader(status)
	if _, err := r.Write([]byte(errorReason)); err != nil {
		return err
	}
	return nil
}

// WriteHeader is overridden to remember the Status Code that has been written.
// Note that using this method, the status value is only written when
// - calling WriteEntity
// - or directly WriteAsXml,WriteAsJson.
// - or if the status is 204 (http.StatusNoContent)
func (r *Response) WriteHeader(httpStatus int) {
	r.statusCode = httpStatus
	// if 204 then WriteEntity will not be called so we need to pass this code
	if http.StatusNoContent == httpStatus {
		r.ResponseWriter.WriteHeader(httpStatus)
	}
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
// Write is part of http.ResponseWriter interface.
func (r *Response) Write(bytes []byte) (int, error) {
	written, err := r.ResponseWriter.Write(bytes)
	r.contentLength += written
	return written, err
}

// ContentLength returns the number of bytes written for the response content.
// Note that this value is only correct if all data is written through the Response using its Write* methods.
// Data written directly using the underlying http.ResponseWriter is not accounted for.
func (r Response) ContentLength() int {
	return r.contentLength
}

// CloseNotify is part of http.CloseNotifier interface
func (r Response) CloseNotify() <-chan bool {
	return r.ResponseWriter.(http.CloseNotifier).CloseNotify()
}
