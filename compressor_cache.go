package restful

// Copyright 2015 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

import (
	"compress/gzip"
	"compress/zlib"
)

type BoundedCachedCompressors struct {
	gzipWriters     chan *gzip.Writer
	gzipReaders     chan *gzip.Reader
	zlibWriters     chan *zlib.Writer
	writersCapacity int
	readersCapacity int
}

func NewBoundedCachedCompressors(writersCapacity, readersCapacity int) *BoundedCachedCompressors {
	return &BoundedCachedCompressors{
		gzipWriters:     make(chan *gzip.Writer, writersCapacity),
		gzipReaders:     make(chan *gzip.Reader, readersCapacity),
		zlibWriters:     make(chan *zlib.Writer, writersCapacity),
		writersCapacity: writersCapacity,
		readersCapacity: readersCapacity,
	}
}

func (b *BoundedCachedCompressors) AcquireGzipWriter() *gzip.Writer {
	var writer *gzip.Writer
	select {
	case writer, _ = <-b.gzipWriters:
	default:
		// return a new unmanaged one
		writer = newGzipWriter()
	}
	return writer
}

func (b *BoundedCachedCompressors) ReleaseGzipWriter(w *gzip.Writer) {
	// forget the unmanaged ones
	if len(b.gzipWriters) < b.writersCapacity {
		b.gzipWriters <- w
	}
}

func (b *BoundedCachedCompressors) AcquireGzipReader() *gzip.Reader {
	var reader *gzip.Reader
	select {
	case reader, _ = <-b.gzipReaders:
	default:
		// return a new unmanaged one
		reader = newGzipReader()
	}
	return reader
}

func (b *BoundedCachedCompressors) ReleaseGzipReader(r *gzip.Reader) {
	// forget the unmanaged ones
	if len(b.gzipReaders) < b.readersCapacity {
		b.gzipReaders <- r
	}
}

func (b *BoundedCachedCompressors) AcquireZlibWriter() *zlib.Writer {
	var writer *zlib.Writer
	select {
	case writer, _ = <-b.zlibWriters:
	default:
		// return a new unmanaged one
		writer = newZlibWriter()
	}
	return writer
}

func (b *BoundedCachedCompressors) ReleaseZlibWriter(w *zlib.Writer) {
	// forget the unmanaged ones
	if len(b.zlibWriters) < b.writersCapacity {
		b.zlibWriters <- w
	}
}
