package restful

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
)

// Request is a wrapper for a http Request that provides convenience methods
type Request struct {
	Request        *http.Request
	pathParameters map[string]string
}

// Return the Path parameter value by its name
func (self *Request) PathParameter(name string) string {
	return self.pathParameters[name]
}

// Return the (first) Query parameter value by its name
func (self *Request) QueryParameter(name string) string {
	return self.Request.FormValue(name)
}

// Check the Accept header and read the content into the entityReference
func (self *Request) ReadEntity(entityReference interface{}) error {
	contentType := self.Request.Header.Get(HEADER_ContentType)
	defer self.Request.Body.Close()
	buffer, err := ioutil.ReadAll(self.Request.Body)
	if err == nil && MIME_XML == contentType {
		log.Printf("unmarschalling:%#v", entityReference)
		err = xml.Unmarshal(buffer, entityReference)
	} else {
		if err == nil && MIME_JSON == contentType {
			err = json.Unmarshal(buffer, entityReference)
		}
	}
	return err
}
