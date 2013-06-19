package restful

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
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

// QueryParameter returns the (first) Query parameter value by its name
func (r *Request) QueryParameter(name string) string {
	return r.Request.FormValue(name)
}

// ReadEntity checks the Accept header and reads the content into the entityReference
// Closes the request body.
func (r *Request) ReadEntity(entityReference interface{}) error {
	contentType := r.Request.Header.Get(HEADER_ContentType)
	defer r.Request.Body.Close()
	buffer, err := ioutil.ReadAll(r.Request.Body)
	if err == nil && MIME_XML == contentType {
		err = xml.Unmarshal(buffer, entityReference)
	} else {
		if err == nil && MIME_JSON == contentType {
			err = json.Unmarshal(buffer, entityReference)
		}
	}
	return err
}
