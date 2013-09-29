package restful

// Copyright 2013 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

// Request is a wrapper for a http Request that provides convenience methods
type Request struct {
	Request        *http.Request
	pathParameters map[string]string
}

// PathParameter accesses the Path parameter value by its name
func (r *Request) PathParameter(name string) string {
	return r.pathParameters[name]
}

// PathParameters accesses the Path parameter values
func (r *Request) PathParameters() map[string]string {
	return r.pathParameters
}

// QueryParameter returns the (first) Query parameter value by its name
func (r *Request) QueryParameter(name string) string {
	return r.Request.FormValue(name)
}

// HeaderParameter returns the HTTP Header value of a Header name or empty if missing
func (r *Request) HeaderParameter(name string) string {
	return r.Request.Header.Get(name)
}

// ReadEntity checks the Accept header and reads the content into the entityPointer
func (r *Request) ReadEntity(entityPointer interface{}) error {
	contentType := r.Request.Header.Get(HEADER_ContentType)
	buffer, err := ioutil.ReadAll(r.Request.Body)
	if err != nil {
		return err
	}
	if strings.Contains(contentType, MIME_XML) {
		err = xml.Unmarshal(buffer, entityPointer)
	} else {
		if strings.Contains(contentType, MIME_JSON) {
			err = json.Unmarshal(buffer, entityPointer)
		} else {
			err = errors.New("[restful] Unable to unmarshal content of type:" + contentType)
		}
	}
	return err
}
