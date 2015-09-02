package restful

import (
	"encoding/json"
	"encoding/xml"
	"sync"
)

type EntityReader interface {
	// Read a serialized version of the value from the request.
	// The Request may have a decompressing reader.
	Read(req *Request, v interface{}) error
}

type EntityWriter interface {
	// Write an serialized version of the value on the response.
	// The Response may have a compressing writer.
	Write(req *Request, resp *Response, v interface{}) error
}

type entityJSON struct{}

// Read unmarshalls the value from JSON
func (e entityJSON) Read(req *Request, v interface{}) error {
	decoder := json.NewDecoder(req.Request.Body)
	decoder.UseNumber()
	return decoder.Decode(v)
}

// Write marshalls the value to JSON and set the Content-Type Header.
func (e entityJSON) Write(req *Request, resp *Response, v interface{}) error {
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
		resp.Header().Set(HEADER_ContentType, MIME_JSON)
		_, err = resp.Write(output)
		return err
	}
	// not-so-pretty
	resp.Header().Set(HEADER_ContentType, MIME_JSON)
	return json.NewEncoder(resp).Encode(v)
}

type entityXML struct{}

// Read unmarshalls the value from XML
func (e entityXML) Read(req *Request, v interface{}) error {
	return xml.NewDecoder(req.Request.Body).Decode(v)
}

// Write marshalls the value to JSON and set the Content-Type Header.
func (e entityXML) Write(req *Request, resp *Response, v interface{}) error {
	if v == nil { // do not write a nil representation
		return nil
	}
	if resp.prettyPrint {
		// pretty output must be created and written explicitly
		output, err := xml.MarshalIndent(v, " ", " ")
		if err != nil {
			return err
		}
		resp.Header().Set(HEADER_ContentType, MIME_XML)
		_, err = resp.Write([]byte(xml.Header))
		if err != nil {
			return err
		}
		_, err = resp.Write(output)
		return err
	}
	// not-so-pretty
	resp.Header().Set(HEADER_ContentType, MIME_XML)
	return xml.NewEncoder(resp).Encode(v)
}

var entityRegistry = &entityAccessorRegistry{
	protection: new(sync.RWMutex),
	readers:    map[string]EntityReader{},
	writers:    map[string]EntityWriter{},
}

type entityAccessorRegistry struct {
	protection *sync.RWMutex
	readers    map[string]EntityReader
	writers    map[string]EntityWriter
}

func (r *entityAccessorRegistry) RegisterEntityAccessors(mime string, reader EntityReader, writer EntityWriter) {
	r.protection.Lock()
	defer r.protection.Unlock()
	r.readers[mime] = reader
	r.writers[mime] = writer
}

func (r *entityAccessorRegistry) ReaderAt(mime string) (EntityReader, bool) {
	r.protection.RLock()
	defer r.protection.RUnlock()
	er, ok := r.readers[mime]
	return er, ok
}

func (r *entityAccessorRegistry) WriterAt(mime string) (EntityWriter, bool) {
	r.protection.RLock()
	defer r.protection.RUnlock()
	ew, ok := r.writers[mime]
	return ew, ok
}
