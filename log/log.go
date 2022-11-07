package log

import (
	stdlog "log"
	"os"
)

// StdLogger corresponds to a minimal subset of the interface satisfied by stdlib log.Logger
type StdLogger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
}

var Logger StdLogger

func init() {
	// default Logger
	SetLogger(stdlog.New(os.Stderr, "[restful] ", stdlog.LstdFlags|stdlog.Lshortfile))
}

// SetLogger sets the logger for this package
func SetLogger(customLogger StdLogger) {
	Logger = customLogger
}

// Print delegates to the Logger
func Print(v ...interface{}) {
	Logger.Print(v...)
}

// Printf delegates to the Logger
func Printf(format string, v ...interface{}) {
	Logger.Printf(format, v...)
}

//Println delegates to the Logger
func Println(v ...interface{}) {
	Logger.Println(v...)
}

//Fatal delegates to the Logger when fatal
func Fatal(v ...interface{}) {
	Logger.Fatal(v...)
}

//Fatal delegates to the Logger when fatal
func Fatalf(format string, v ...interface{}) {
	Logger.Fatalf(format, v...)
}

//Fatal delegates to the Logger when fatal
func Fatalln(v ...interface{}) {
	Logger.Fatalln(v...)
}
