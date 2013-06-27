package restful

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

// QueryParameter returns the (first) Query parameter value by its name
func (r *Request) QueryParameter(name string) string {
	return r.Request.FormValue(name)
}

// ReadEntity checks the Accept header and reads the content into the entityPointer
func (r *Request) ReadEntity(entityPointer interface{}) error {
	contentType := r.Request.Header.Get(HEADER_ContentType)
	buffer, err := ioutil.ReadAll(r.Request.Body)
	if err != nil {
		return err
	}
	if strings.Contains(contentType, MIME_XML) {
		err = xml.Unmarshal(buffer, entityReference)
	} else {
		if strings.Contains(contentType, MIME_JSON) {
			err = json.Unmarshal(buffer, entityReference)
		} else {
			err = errors.New("[restful] Unable to unmarshal content of type:" + contentType)
		}
	}
	return err
}
