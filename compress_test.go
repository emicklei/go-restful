package restful

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// go test -v -test.run TestGzip ...restful
func TestGzip(t *testing.T) {
	EnableContentEncoding = true
	httpRequest, _ := http.NewRequest("GET", "/test", nil)
	httpRequest.Header.Set("Accept-Encoding", "gzip,deflate")
	httpWriter := httptest.NewRecorder()
	wanted, encoding := wantsCompressedResponse(httpRequest)
	if !wanted {
		t.Fatal("should accept gzip")
	}
	if encoding != "gzip" {
		t.Fatal("expected gzip")
	}
	c, err := NewCompressingResponseWriter(httpWriter, encoding)
	if err != nil {
		t.Fatal(err.Error())
	}
	c.Write([]byte("Hello World"))
	c.Close()
	if httpWriter.Header().Get("Content-Encoding") != "gzip" {
		t.Fatal("Missing gzip header")
	}
	reader, err := gzip.NewReader(httpWriter.Body)
	if err != nil {
		t.Fatal(err.Error())
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatal(err.Error())
	}
	if got, want := string(data), "Hello World"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestDeflate(t *testing.T) {
	EnableContentEncoding = true
	httpRequest, _ := http.NewRequest("GET", "/test", nil)
	httpRequest.Header.Set("Accept-Encoding", "deflate,gzip")
	httpWriter := httptest.NewRecorder()
	wanted, encoding := wantsCompressedResponse(httpRequest)
	if !wanted {
		t.Fatal("should accept deflate")
	}
	if encoding != "deflate" {
		t.Fatal("expected deflate")
	}
	c, err := NewCompressingResponseWriter(httpWriter, encoding)
	if err != nil {
		t.Fatal(err.Error())
	}
	c.Write([]byte("Hello World"))
	c.Close()
	if httpWriter.Header().Get("Content-Encoding") != "deflate" {
		t.Fatal("Missing deflate header")
	}
	reader, err := zlib.NewReader(httpWriter.Body)
	if err != nil {
		t.Fatal(err.Error())
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatal(err.Error())
	}
	if got, want := string(data), "Hello World"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestGzipDecompressRequestBody(t *testing.T) {
	b := new(bytes.Buffer)
	w := newGzipWriter()
	w.Reset(b)
	io.WriteString(w, "hello") // -> [31 139 8 0 0 9 110 136 4 255]
	w.Flush()
	w.Close()

	gzipReader := gzipReaderPool.Get().(*gzip.Reader)
	gzipReader.Reset(bytes.NewReader(b.Bytes()))
	s, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		t.Error(err)
	}
	if got, want := string(s), "hello"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
