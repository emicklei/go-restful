package restful

// Copyright 2015 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

import (
	"compress/gzip"
	"compress/zlib"
)

type CompressorProvider interface {
	// Returns a *gzip.Writer which needs to be released later.
	// Before using it, call Reset().
	AcquireGzipWriter() *gzip.Writer

	// Releases an aqcuired *gzip.Writer.
	ReleaseGzipWriter(w *gzip.Writer)

	// Returns a *gzip.Reader which needs to be released later.
	AcquireGzipReader() *gzip.Reader

	// Releases an aqcuired *gzip.Reader.
	ReleaseGzipReader(w *gzip.Reader)

	// Returns a *zlib.Writer which needs to be released later.
	// Before using it, call Reset().
	AcquireZlibWriter() *zlib.Writer

	// Releases an aqcuired *zlib.Writer.
	ReleaseZlibWriter(w *zlib.Writer)
}

// DefaultCompressorProvider is the actual provider of compressors (zlib or gzip).
var DefaultCompressorProvider CompressorProvider

func init() {
	DefaultCompressorProvider = NewSyncPoolCompessors()
}
