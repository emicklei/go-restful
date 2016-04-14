package restful

import (
	"io"
	"io/ioutil"

	"gopkg.in/vmihailenco/msgpack.v2"
)

// NewEntityAccessorMPack returns a new EntityReaderWriter for accessing MessagePack content.
// This package is not initialized with such an accessor using the MIME_MSGPACK contentType.
func NewEntityAccessorMsgPack(contentType string) EntityReaderWriter {
	return entityMsgPackAccess{ContentType: contentType}
}

// entityOctetAccess is a EntityReaderWriter for Octet encoding
type entityMsgPackAccess struct {
	// This is used for setting the Content-Type header when writing
	ContentType string
}

// Read unmarshalls the value from byte slice and using msgpack to unmarshal
func (e entityMsgPackAccess) Read(req *Request, v interface{}) error {
	data, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(data, v)
}

// Write marshals the value to byte slice and set the Content-Type Header.
func (e entityMsgPackAccess) Write(resp *Response, status int, v interface{}) error {
	return writeMsgPack(resp, status, e.ContentType, v)
}

// writeMsgPack marshals the value to byte slice and set the Content-Type Header.
func writeMsgPack(resp *Response, status int, contentType string, v interface{}) error {
	if v == nil {
		resp.WriteHeader(status)
		// do not write a nil representation
		return nil
	}
	resp.Header().Set(HEADER_ContentType, contentType)
	resp.WriteHeader(status)

	m, err := msgpack.Marshal(v)
	if err != nil {
		return err
	}
	_, err = io.WriteString(resp, string(m[:]))
	return err
}
