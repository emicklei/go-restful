package restful

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"strings"
	"sync"
)

var (
	// entities are the registered entities
	entities map[string]EntityEncoder
	// m protects the entities map
	entitiesM sync.RWMutex
)

// initialize the entities map and register the default entity encoders
func init() {
	entities = make(map[string]EntityEncoder)
	RegisterEntityEncoder(&JSONEntity{})
	RegisterEntityEncoder(&XMLEntity{})
}

// EntityEncoder describes how an entity should be encoded/decoded in ReadEntity/WriteEntity
//
// It can receive a Request, for additional interaction with headers
// Because of the stateful nature of setting a request, an EntityEncoder must be
// able to return a new instance of itself for each call to New()
type EntityEncoder interface {
	// New must return a new instance of the encoder if it is required
	// to hold state through SetRequest
	New() EntityEncoder
	// MIME must return the MIME type to be used in the Content-Type header
	MIME() string
	// SetRequest will receive the current request
	// It must only operate on instances created with New
	SetRequest(*Request)
	//  SetResponse will receive the current response
	// It must only operate on instances created with New
	SetResponse(*Response)
	// Marshal to be used in Response writing
	Marshal(v interface{}) ([]byte, error)
	// MarshalIndent for PrettyPrint cases
	MarshalIndent(v interface{}, prefix, indent string) ([]byte, error)
	// Unmarshal to be used in Request reading
	Unmarshal(b []byte, entityPointer interface{}) error
}

// RegisterEntity registers an entity type to be handled by Requests/Responses
func RegisterEntityEncoder(e EntityEncoder) {
	entitiesM.Lock()
	entities[e.MIME()] = e
	entitiesM.Unlock()
}

// EntityForMIME returns an EntityEncoder for the given MIME type or nil if not registered
func EntityEncoderForMIME(mime string) EntityEncoder {
	entitiesM.RLock()
	if e, ok := entities[mime]; ok {
		entitiesM.RUnlock()
		return e.New()
	} else {
		entitiesM.RUnlock()
		return nil
	}
}

// EntityForContentType returns an EntityEncoder for the given Content-Type header or nil if not registered
func EntityEncoderForContentType(contentType string) EntityEncoder {
	entitiesM.RLock()
	for mime, e := range entities {
		if strings.Contains(contentType, mime) {
			entitiesM.RUnlock()
			return e.New()
		}
	}
	return nil
}

// JSONEntity describes the JSON entity encoding
type JSONEntity struct{}

// New implementing the EntityEncoder interface
func (e *JSONEntity) New() EntityEncoder {
	// non-stateful, no need to allocate
	return e
}

// MIME will return MIME_JSON
func (e *JSONEntity) MIME() string {
	return MIME_JSON
}

// SetRequest to set the request
// no-op
func (e *JSONEntity) SetRequest(r *Request) {}

// SetResponse is a no-op
func (e *JSONEntity) SetResponse(r *Response) {}

// Marshal passthrough json.Marshal
func (e *JSONEntity) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// MarshalIndent passthrough json.MarshalIndent
func (e *JSONEntity) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

// Unmarshal passthrough json.Unmarshal
func (e *JSONEntity) Unmarshal(b []byte, entityPointer interface{}) error {
	decoder := json.NewDecoder(bytes.NewReader(b))
	decoder.UseNumber()
	return decoder.Decode(entityPointer)
}

// XMLEntity describes the XML entity encoding
type XMLEntity struct{}

// New implementing the EntityEncoder interface
func (e *XMLEntity) New() EntityEncoder {
	// non-stateful
	return e
}

// MIME will return MIME_XMl
func (e *XMLEntity) MIME() string {
	return MIME_XML
}

// SetRequest for current state: no-op
func (e *XMLEntity) SetRequest(r *Request) {}

// SetResponse is a no-op
func (e *XMLEntity) SetResponse(r *Response) {}

// withHeader will include the xml.Header in the output
func (e *XMLEntity) withHeader(xmlBytes []byte, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	b := []byte(xml.Header)
	b = append(b, xmlBytes...)
	return b, nil
}

// Marshal will encode the value to XMl
func (e *XMLEntity) Marshal(v interface{}) ([]byte, error) {
	return e.withHeader(xml.Marshal(v))
}

// MarshalIndent will encode the vaue to XML for pretty printing
func (e *XMLEntity) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return e.withHeader(xml.MarshalIndent(v, prefix, indent))
}

// Unmarshal to decode the XML
func (e *XMLEntity) Unmarshal(b []byte, entityPointer interface{}) error {
	return xml.Unmarshal(b, entityPointer)
}
