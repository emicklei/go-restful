package restful

// Copyright 2014 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.
import (
	stdlog "log"
	"os"

	"github.com/go-log/log"
	"github.com/go-log/log/print"
)

var (
	// Logger is the package-wide logger.
	Logger log.Logger

	trace bool = false
	traceLogger log.Logger
)

func init() {
	Logger = print.New(stdlog.New(os.Stderr, "[restful] ", stdlog.LstdFlags|stdlog.Lshortfile))
	traceLogger = Logger
}

// TraceLogger enables detailed logging of HTTP request matching and filter invocation. Default no logger is set.
// You may call EnableTracing() directly to enable trace logging to the package-wide logger.
//
// Deprecated: Please use TraceLogLogger(print.New(logger)) instead, with New from github.com/go-log/log/print.
func TraceLogger(logger print.Printer) {
	TraceLogLogger(print.New(logger))
}

// TraceLogLogger enables detailed logging of HTTP request matching and filter invocation. Default no logger is set.
// You may call EnableTracing() directly to enable trace logging to the package-wide logger.
func TraceLogLogger(logger log.Logger) {
	traceLogger = logger
	EnableTracing(logger != nil)
}

// SetLogger sets the package-wide logger.
//
// Deprecated: Please set Logger instead.  If you need to convert from
// an implementation that provides Print and Printf, use
// github.com/go-log/log/print's New().
func SetLogger(customLogger print.Printer) {
	Logger = print.New(customLogger)
}

// EnableTracing can be used to Trace logging on and off.
func EnableTracing(enabled bool) {
	trace = enabled
}
