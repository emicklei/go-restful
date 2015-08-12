package restful

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"sync"
)

// gzipWriterPool is used to get reusable zippers.
var gzipWriterPool = &sync.Pool{
	New: func() interface{} {
		return newGzipWriter()
	},
}

func newGzipWriter() *gzip.Writer {
	// create with an empty bytes writer; it will be replaced before using the gzipWriter
	writer, err := gzip.NewWriterLevel(new(bytes.Buffer), gzip.BestSpeed)
	if err != nil {
		panic(err.Error())
	}
	return writer
}

// gzipReaderPool is used to get reusable zippers.
var gzipReaderPool = &sync.Pool{
	New: func() interface{} {
		return newGzipReader()
	},
}

func newGzipReader() *gzip.Reader {
	// create with an empty reader (but with GZIP header); it will be replaced before using the gzipReader
	w := gzipWriterPool.Get().(*gzip.Writer)
	b := new(bytes.Buffer)
	w.Reset(b)
	w.Flush()
	w.Close()
	reader, err := gzip.NewReader(bytes.NewReader(b.Bytes()))
	if err != nil {
		panic(err.Error())
	}
	return reader
}

// zlibWriterPool is used to get reusable zippers.
var zlibWriterPool = &sync.Pool{
	New: func() interface{} {
		return newZlibWriter()
	},
}

func newZlibWriter() *zlib.Writer {
	writer, err := zlib.NewWriterLevel(new(bytes.Buffer), gzip.BestSpeed)
	if err != nil {
		panic(err.Error())
	}
	return writer
}
