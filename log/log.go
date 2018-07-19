// Package log is deprecated.  Please use restful.SetLogLogger instead.
package log

import (
	"github.com/emicklei/go-restful"
	"github.com/go-log/log/print"
)

// StdLogger corresponds to a minimal subset of the interface satisfied by stdlib log.Logger
type StdLogger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
}

// SetLogger sets the logger for this package
//
// Deprecated: Please set restful.Logger instead.
func SetLogger(customLogger StdLogger) {
	restful.Logger = print.New(customLogger)
}

// Print delegates to the Logger
//
// Deprecated: Please use restful.Logger.Log instead.
func Print(v ...interface{}) {
	restful.Logger.Log(v...)
}

// Printf delegates to the Logger
//
// Deprecated: Please use restful.Logger.Logf instead.
func Printf(format string, v ...interface{}) {
	restful.Logger.Logf(format, v...)
}
