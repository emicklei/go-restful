// Copyright (c) 2013 Ernest Micklei. All rights reserved.

package restful

import (
	"compress/gzip"
	"compress/zlib"
	"errors"
	"io"
	"net/http"
	"strings"
)

// OBSOLETE : use restful.DefaultContainer.EnableContentEncoding(true) to change this setting.
var EnableContentEncoding = false

// CompressingResponseWriter is a http.ResponseWriter that can perform content encoding (gzip and zlib)
type CompressingResponseWriter struct {
	writer     http.ResponseWriter
	compressor io.WriteCloser
}

// Header is part of http.ResponseWriter interface
func (c *CompressingResponseWriter) Header() http.Header {
	return c.writer.Header()
}

// WriteHeader is part of http.ResponseWriter interface
func (c *CompressingResponseWriter) WriteHeader(status int) {
	c.writer.WriteHeader(status)
}

// Write is part of http.ResponseWriter interface
// It is passed through the compressor
func (c *CompressingResponseWriter) Write(bytes []byte) (int, error) {
	return c.compressor.Write(bytes)
}

// Close the underlying compressor
func (c *CompressingResponseWriter) Close() {
	c.compressor.Close()
}

// WantsCompressedResponse reads the Accept-Encoding header to see if and which encoding is requested.
func wantsCompressedResponse(httpRequest *http.Request) (bool, string) {
	header := httpRequest.Header.Get(HEADER_AcceptEncoding)
	gi := strings.Index(header, ENCODING_GZIP)
	zi := strings.Index(header, ENCODING_DEFLATE)
	// use in order of appearance
	if gi == -1 {
		return zi != -1, ENCODING_DEFLATE
	} else if zi == -1 {
		return gi != -1, ENCODING_GZIP
	} else {
		if gi < zi {
			return true, ENCODING_GZIP
		} else {
			return true, ENCODING_DEFLATE
		}
	}
}

// NewCompressingResponseWriter create a CompressingResponseWriter for a known encoding = {gzip,deflate}
func NewCompressingResponseWriter(httpWriter http.ResponseWriter, encoding string) (*CompressingResponseWriter, error) {
	httpWriter.Header().Set(HEADER_ContentEncoding, encoding)
	c := new(CompressingResponseWriter)
	c.writer = httpWriter
	var err error
	if ENCODING_GZIP == encoding {
		c.compressor, err = gzip.NewWriterLevel(httpWriter, gzip.BestSpeed)
		if err != nil {
			return nil, err
		}
	} else if ENCODING_DEFLATE == encoding {
		c.compressor, err = zlib.NewWriterLevel(httpWriter, zlib.BestSpeed)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Unknown encoding:" + encoding)
	}
	return c, err
}
