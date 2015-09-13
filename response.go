package restful

// Copyright 2013 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

import (
	"errors"
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
	prettyPrint   bool     // controls the indentation feature of XML and JSON serialization. It is initialized using var PrettyPrintResponses.
	err           error    // err property is kept when WriteError is called
}

// Creates a new response based on a http ResponseWriter.
func NewResponse(httpWriter http.ResponseWriter) *Response {
	return &Response{httpWriter, "", []string{}, http.StatusOK, 0, PrettyPrintResponses, nil} // empty content-types
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

// PrettyPrint changes whether this response must produce pretty (line-by-line, indented) JSON or XML output.
func (r *Response) PrettyPrint(bePretty bool) {
	r.prettyPrint = bePretty
}

// AddHeader is a shortcut for .Header().Add(header,value)
func (r Response) AddHeader(header string, value string) Response {
	r.Header().Add(header, value)
	return r
}

// SetRequestAccepts tells the response what Mime-type(s) the HTTP request said it wants to accept. Exposed for testing.
func (r *Response) SetRequestAccepts(mime string) {
	r.requestAccept = mime
}

// EntityWriter returns the registered EntityWriter that the entity (requested resource)
// can write according to what the request wants (Accept) and what the Route can produce or what the restful defaults say.
// If called before WriteEntity and WriteHeader then a false return value can be used to write a 406: Not Acceptable.
func (r *Response) EntityWriter() (EntityReaderWriter, bool) {
	for _, qualifiedMime := range strings.Split(r.requestAccept, ",") {
		mime := strings.Trim(strings.Split(qualifiedMime, ";")[0], " ")
		if 0 == len(mime) || mime == "*/*" {
			for _, each := range r.routeProduces {
				if MIME_JSON == each {
					return entityAccessRegistry.AccessorAt(MIME_JSON)
				}
				if MIME_XML == each {
					return entityAccessRegistry.AccessorAt(MIME_XML)
				}
			}
		} else { // mime is not blank; see if we have a match in Produces
			for _, each := range r.routeProduces {
				if mime == each {
					if MIME_JSON == each {
						return entityAccessRegistry.AccessorAt(MIME_JSON)
					}
					if MIME_XML == each {
						return entityAccessRegistry.AccessorAt(MIME_XML)
					}
				}
			}
		}
	}
	writer, ok := entityAccessRegistry.AccessorAt(r.requestAccept)
	if !ok {
		// if not registered then fallback to the defaults (if set)
		if DefaultResponseMimeType == MIME_JSON {
			return entityAccessRegistry.AccessorAt(MIME_JSON)
		}
		if DefaultResponseMimeType == MIME_XML {
			return entityAccessRegistry.AccessorAt(MIME_XML)
		}
	}
	return writer, ok
}

// WriteEntity marshals the value using the representation denoted by the Accept Header and the registered EntityWriters.
// If no Accept header is specified (or */*) then return the Content-Type as specified by the first in the Route.Produces.
// If an Accept header is specified then return the Content-Type as specified by the first in the Route.Produces that is matched with the Accept header.
// If the value is nil then nothing is written. You may want to call WriteHeader(http.StatusNotFound) instead.
// Current implementation ignores any q-parameters in the Accept Header.
func (r *Response) WriteEntity(value interface{}) error {
	if value == nil { // do not write a nil representation
		return nil
	}
	writer, ok := r.EntityWriter()
	if !ok {
		return errors.New("response cannot be written in an acceptable content-type.")
	}
	return writer.Write(r, value)
}

// WriteAsXml is a convenience method for writing a value in xml (requires Xml tags on the value)
// It uses the standard encoding/xml package for marshalling the valuel ; not using a registered EntityReaderWriter.
func (r *Response) WriteAsXml(value interface{}) error {
	return writeXML(r, MIME_XML, value)
}

// WriteAsJson is a convenience method for writing a value in json.
// It uses the standard encoding/json package for marshalling the valuel ; not using a registered EntityReaderWriter.
func (r *Response) WriteAsJson(value interface{}) error {
	return writeJSON(r, MIME_JSON, value)
}

// WriteJson is a convenience method for writing a value in Json with a given Content-Type.
// It uses the standard encoding/json package for marshalling the valuel ; not using a registered EntityReaderWriter.
func (r *Response) WriteJson(value interface{}, contentType string) error {
	return writeJSON(r, contentType, value)
}

// WriteError write the http status and the error string on the response.
func (r *Response) WriteError(httpStatus int, err error) error {
	r.err = err
	return r.WriteErrorString(httpStatus, err.Error())
}

// WriteServiceError is a convenience method for a responding with a status and a ServiceError
func (r *Response) WriteServiceError(httpStatus int, err ServiceError) error {
	r.WriteHeader(httpStatus)
	return r.WriteEntity(err)
}

// WriteErrorString is a convenience method for an error status with the actual error
func (r *Response) WriteErrorString(httpStatus int, errorReason string) error {
	r.WriteHeader(httpStatus)
	if _, err := r.Write([]byte(errorReason)); err != nil {
		return err
	}
	return nil
}

// WriteHeader is overridden to remember the Status Code that has been written.
func (r *Response) WriteHeader(httpStatus int) {
	r.statusCode = httpStatus
	r.ResponseWriter.WriteHeader(httpStatus)
}

// StatusCode returns the code that has been written using WriteHeader.
// If WriteHeader, WriteEntity or WriteAsXml has not been called (yet) then return 200 OK.
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

// Error returns the err created by WriteError
func (r Response) Error() error {
	return r.err
}
