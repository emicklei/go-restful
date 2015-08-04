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
		return newGzipper()
	},
}

func newGzipper() *gzip.Writer {
	writer, err := gzip.NewWriterLevel(new(bytes.Buffer), gzip.BestSpeed)
	if err != nil {
		panic(err.Error())
	}
	return writer
}

// zlibWriterPool is used to get reusable zippers.
var zlibWriterPool = &sync.Pool{
	New: func() interface{} {
		return newZlibber()
	},
}

func newZlibber() *zlib.Writer {
	writer, err := zlib.NewWriterLevel(new(bytes.Buffer), gzip.BestSpeed)
	if err != nil {
		panic(err.Error())
	}
	return writer
}
