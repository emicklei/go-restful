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
	"time"
)

type flushSnapshotRecorder struct {
	*httptest.ResponseRecorder
	snapshots [][]byte
}

func (r *flushSnapshotRecorder) Flush() {
	r.snapshots = append(r.snapshots, append([]byte(nil), r.Body.Bytes()...))
	r.ResponseRecorder.Flush()
}

func compressedStreamReader(encoding string, input io.Reader) (io.ReadCloser, error) {
	if encoding == ENCODING_GZIP {
		return gzip.NewReader(input)
	}
	return zlib.NewReader(input)
}

func TestCompressingResponseWriterFlushesPendingData(t *testing.T) {
	for _, encoding := range []string{ENCODING_GZIP, ENCODING_DEFLATE} {
		t.Run(encoding, func(t *testing.T) {
			recorder := &flushSnapshotRecorder{ResponseRecorder: httptest.NewRecorder()}
			writer, err := NewCompressingResponseWriter(recorder, encoding)
			if err != nil {
				t.Fatal(err)
			}
			defer writer.Close()
			var want []byte
			for _, chunk := range []string{"first small chunk", " and a second chunk"} {
				want = append(want, chunk...)
				if _, err := io.WriteString(writer, chunk); err != nil {
					t.Fatal(err)
				}
				beforeFlushes := len(recorder.snapshots)
				writer.Flush()
				if len(recorder.snapshots) != beforeFlushes+1 {
					t.Fatal("Flush did not call the underlying flusher exactly once")
				}
				snapshot := recorder.snapshots[len(recorder.snapshots)-1]
				reader, err := compressedStreamReader(encoding, bytes.NewReader(snapshot))
				if err != nil {
					t.Fatal(err)
				}
				got := make([]byte, len(want))
				_, err = io.ReadFull(reader, got)
				reader.Close()
				if err != nil {
					t.Fatalf("flushed stream did not contain pending data: %v", err)
				}
				if !bytes.Equal(got, want) {
					t.Fatalf("flushed %q, want %q", got, want)
				}
			}
		})
	}
}

func TestContainerStreamsCompressedResponseBeforeClose(t *testing.T) {
	for _, encoding := range []string{ENCODING_GZIP, ENCODING_DEFLATE} {
		t.Run(encoding, func(t *testing.T) {
			allowNext := make(chan struct{})
			service := new(WebService)
			service.Route(service.GET("/stream").To(func(req *Request, resp *Response) {
				if _, err := io.WriteString(resp, "first chunk"); err != nil {
					return
				}
				resp.Flush()
				// Keep the compressor open until the client receives the first
				// chunk. A close must not be what makes that chunk available.
				select {
				case <-allowNext:
				case <-req.Request.Context().Done():
					return
				}
				io.WriteString(resp, "second chunk")
				resp.Flush()
			}))
			container := NewContainer()
			container.EnableContentEncoding(true)
			container.Add(service)
			server := httptest.NewServer(container)
			defer server.Close()
			client := server.Client()
			client.Timeout = 3 * time.Second
			req, err := http.NewRequest(http.MethodGet, server.URL+"/stream", nil)
			if err != nil {
				t.Fatal(err)
			}
			// Supplying this explicitly prevents the Go client from adding
			// gzip itself and transparently decoding the response.
			req.Header.Set(HEADER_AcceptEncoding, encoding)
			response, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer response.Body.Close()
			if got := response.Header.Get(HEADER_ContentEncoding); got != encoding {
				t.Fatalf("Content-Encoding = %q, want %q", got, encoding)
			}
			reader, err := compressedStreamReader(encoding, response.Body)
			if err != nil {
				t.Fatal(err)
			}
			defer reader.Close()
			first := make([]byte, len("first chunk"))
			if _, err := io.ReadFull(reader, first); err != nil {
				t.Fatalf("first chunk unavailable before close: %v", err)
			}
			if string(first) != "first chunk" {
				t.Fatalf("received %q, want first chunk", first)
			}
			close(allowNext)
			rest, err := ioutil.ReadAll(reader)
			if err != nil {
				t.Fatal(err)
			}
			if string(rest) != "second chunk" {
				t.Fatalf("received %q, want second chunk", rest)
			}
		})
	}
}

type nonFlushingResponseWriter struct {
	header http.Header
	body   bytes.Buffer
}

func (w *nonFlushingResponseWriter) Header() http.Header { return w.header }
func (w *nonFlushingResponseWriter) WriteHeader(int)     {}
func (w *nonFlushingResponseWriter) Write(p []byte) (int, error) {
	return w.body.Write(p)
}

func TestCompressingResponseWriterFlushWithoutUnderlyingFlusher(t *testing.T) {
	for _, encoding := range []string{ENCODING_GZIP, ENCODING_DEFLATE} {
		t.Run(encoding, func(t *testing.T) {
			recorder := &nonFlushingResponseWriter{header: make(http.Header)}
			writer, err := NewCompressingResponseWriter(recorder, encoding)
			if err != nil {
				t.Fatal(err)
			}
			defer writer.Close()
			if _, err := io.WriteString(writer, "small buffered chunk"); err != nil {
				t.Fatal(err)
			}
			before := append([]byte(nil), recorder.body.Bytes()...)
			writer.Flush()
			if !bytes.Equal(recorder.body.Bytes(), before) {
				t.Fatal("Flush changed output without an underlying flusher")
			}
		})
	}
}

func TestCompressingResponseWriterFlushAfterClose(t *testing.T) {
	for _, encoding := range []string{ENCODING_GZIP, ENCODING_DEFLATE} {
		t.Run(encoding, func(t *testing.T) {
			recorder := &flushSnapshotRecorder{ResponseRecorder: httptest.NewRecorder()}
			writer, err := NewCompressingResponseWriter(recorder, encoding)
			if err != nil {
				t.Fatal(err)
			}
			if err := writer.Close(); err != nil {
				t.Fatal(err)
			}
			before := append([]byte(nil), recorder.Body.Bytes()...)
			writer.Flush()
			if len(recorder.snapshots) != 1 || !bytes.Equal(before, recorder.snapshots[0]) {
				t.Fatal("Flush after Close did not preserve the completed output")
			}
		})
	}
}
