package log

import (
	stdlog "log"
	"os"
)

// Logger corresponds to a subset of the interface satisfied by stdlib log.Logger
type StdLogger interface {
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

var Logger StdLogger

func init() {
	// default Logger
	SetLogger(stdlog.New(os.Stdout, "[restful] ", stdlog.LstdFlags|stdlog.Lshortfile))
}

func SetLogger(customLogger StdLogger) {
	Logger = customLogger
}

func Fatal(v ...interface{}) {
	Logger.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	Logger.Fatalf(format, v...)
}

func Fatalln(v ...interface{}) {
	Logger.Fatalln(v...)
}

func Print(v ...interface{}) {
	Logger.Print(v...)
}

func Printf(format string, v ...interface{}) {
	Logger.Printf(format, v...)
}

func Println(v ...interface{}) {
	Logger.Println(v...)
}
