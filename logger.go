package restful

// Copyright 2014 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

var trace bool = false
var traceLogger Logger

// Logger corresponds to a subset of the interface satisfied by stdlib log.Logger
type Logger interface {
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
	// Flags() int
	// Output(calldepth int, s string) error
	// Panic(v ...interface{})
	// Panicf(format string, v ...interface{})
	// Panicln(v ...interface{})
	// Prefix() string
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	// SetFlags(flag int)
	// SetPrefix(prefix string)
}

// TraceLogger enables detailed logging of Http request matching and filter invocation. Default no logger is set.
func TraceLogger(logger Logger) {
	traceLogger = logger
	trace = logger != nil
}
