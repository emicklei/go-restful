package restful

import (
	"encoding/json"
	"encoding/xml"
	"strings"
	"sync"
)

type EntityReaderWriter interface {
	// Read a serialized version of the value from the request.
	// The Request may have a decompressing reader. Depends on Content-Encoding.
	Read(req *Request, v interface{}) error

	// Write an serialized version of the value on the response.
	// The Response may have a compressing writer. Depends on Accept-Encoding.
	Write(resp *Response, v interface{}) error
}

var entityAccessRegistry = &entityReaderWriters{
	protection: new(sync.RWMutex),
	accessors:  map[string]EntityReaderWriter{},
}

type entityReaderWriters struct {
	protection *sync.RWMutex
	accessors  map[string]EntityReaderWriter
}

func init() {
	jsonRW := JSONEntityCodec{ContentType: MIME_JSON}
	xmlRW := XMLEntityCodec{ContentType: MIME_XML}
	RegisterEntityAccessor(MIME_JSON, jsonRW)
	RegisterEntityAccessor(MIME_XML, xmlRW)
}

func RegisterEntityAccessor(mime string, erw EntityReaderWriter) {
	entityAccessRegistry.protection.Lock()
	defer entityAccessRegistry.protection.Unlock()
	entityAccessRegistry.accessors[mime] = erw
}

func (r *entityReaderWriters) AccessorAt(mime string) (EntityReaderWriter, bool) {
	r.protection.RLock()
	defer r.protection.RUnlock()
	er, ok := r.accessors[mime]
	if !ok {
		// retry with reverse lookup
		// more expensive but we are in an exceptional situation anyway
		for k, v := range r.accessors {
			if strings.Contains(mime, k) {
				return v, true
			}
		}
	}
	return er, ok
}

type XMLEntityCodec struct {
	// This is used for setting the Content-Type header when writing
	ContentType string
}

// Read unmarshalls the value from XML
func (e XMLEntityCodec) Read(req *Request, v interface{}) error {
	return xml.NewDecoder(req.Request.Body).Decode(v)
}

// Write marshalls the value to JSON and set the Content-Type Header.
func (e XMLEntityCodec) Write(resp *Response, v interface{}) error {
	if v == nil { // do not write a nil representation
		return nil
	}
	if resp.prettyPrint {
		// pretty output must be created and written explicitly
		output, err := xml.MarshalIndent(v, " ", " ")
		if err != nil {
			return err
		}
		resp.Header().Set(HEADER_ContentType, e.ContentType)
		_, err = resp.Write([]byte(xml.Header))
		if err != nil {
			return err
		}
		_, err = resp.Write(output)
		return err
	}
	// not-so-pretty
	resp.Header().Set(HEADER_ContentType, e.ContentType)
	return xml.NewEncoder(resp).Encode(v)
}

type JSONEntityCodec struct {
	// This is used for setting the Content-Type header when writing
	ContentType string
}

// Read unmarshalls the value from JSON
func (e JSONEntityCodec) Read(req *Request, v interface{}) error {
	decoder := json.NewDecoder(req.Request.Body)
	decoder.UseNumber()
	return decoder.Decode(v)
}

// Write marshalls the value to JSON and set the Content-Type Header.
func (e JSONEntityCodec) Write(resp *Response, v interface{}) error {
	if v == nil {
		// do not write a nil representation
		return nil
	}
	if resp.prettyPrint {
		// pretty output must be created and written explicitly
		output, err := json.MarshalIndent(v, " ", " ")
		if err != nil {
			return err
		}
		resp.Header().Set(HEADER_ContentType, e.ContentType)
		_, err = resp.Write(output)
		return err
	}
	// not-so-pretty
	resp.Header().Set(HEADER_ContentType, e.ContentType)
	return json.NewEncoder(resp).Encode(v)
}
