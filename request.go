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
func (self *Request) PathParameter(name string) string {
	return self.pathParameters[name]
}

// QueryParameter returns the (first) Query parameter value by its name
func (self *Request) QueryParameter(name string) string {
	return self.Request.FormValue(name)
}

// ReadEntity check the Accept header and reads the content into the entityReference
func (self *Request) ReadEntity(entityReference interface{}) error {
	contentType := self.Request.Header.Get(HEADER_ContentType)
	defer self.Request.Body.Close()
	buffer, err := ioutil.ReadAll(self.Request.Body)
	if err == nil && MIME_XML == contentType {
		err = xml.Unmarshal(buffer, entityReference)
	} else {
		if err == nil && MIME_JSON == contentType {
			err = json.Unmarshal(buffer, entityReference)
		}
	}
	return err
}
