//  log  -  Simple logging.
//
//  E.B.Smith  -  November, 2014


package log


import (
    "io"
    "os"
    "fmt"
    "syscall"
    "runtime"
    "path"
    )


type LogLevelType int
const (
    LevelDebug LogLevelType = iota
    LevelStart
    LevelExit
    LevelInfo
    LevelWarning
    LevelError
    LevelAll = LevelDebug
    )


var LogLevel    LogLevelType    = LevelDebug
var logWriter   io.WriteCloser  = os.Stderr


func SetFilename(filename string) {
    if logWriter.Close != nil { logWriter.Close() }
    var error error
    var flags int = syscall.O_APPEND | syscall.O_CREAT | syscall.O_WRONLY
    var mode os.FileMode = os.ModeAppend | os.ModePerm
    logWriter, error = os.OpenFile(filename, flags, mode)
    if error != nil {
        logWriter = os.Stderr
        Error("Error: Can't open log file '%s' for reading: %v.", filename, error)
    }
}


func StackWithError(error interface{}) {
    trace := make([]byte, 1024)
    count := runtime.Stack(trace, true)
    fmt.Fprintf(logWriter, "Panic!! '%v'.\n", error)
    fmt.Fprintf(logWriter, "Stack of %d bytes: %s\n", count, trace)
}


func PrettyStackString(stackLevel int) string {
    _, filename, linenumber, _ := runtime.Caller(stackLevel)
    filename = path.Base(filename)
    i := len(filename)
    if i > 26 {
        i = 26
    }
    return fmt.Sprintf("%s:%d", filename[:i], linenumber)
}


func LogFunctionName() {
    pc, filename, linenumber, _ := runtime.Caller(1)
    filename = path.Base(filename)
    i := len(filename)
    if i > 26 {
        i = 26
    }
    message := fmt.Sprintf("function %s.", runtime.FuncForPC(pc).Name())
    fmt.Fprintf(logWriter, "%26s:%-4d %s %s\n", filename[:i], linenumber, " Info", message)
}


func logRaw(logLevel LogLevelType, format string, args ...interface{}) {

    LevelNames := []string {
        "Debug",
        "Start",
        " Exit",
        " Info",
        " Warn",
        "Error",
    }

    if logLevel < LogLevel { return }
    if logLevel < LevelDebug || logLevel > LevelError { logLevel = LevelError }

    _, filename, linenumber, _ := runtime.Caller(2)
    filename = path.Base(filename)
    i := len(filename)
    if i > 26 {
        i = 26
    }

    var message = fmt.Sprintf(format, args...)
    fmt.Fprintf(logWriter, "%26s:%-4d %s %s\n", filename[:i], linenumber, LevelNames[logLevel], message)
}


func Debug(format string, args ...interface{})          { logRaw(LevelDebug, format, args...) }
func Start(format string, args ...interface{})          { logRaw(LevelStart, format, args...) }
func Exit(format string, args ...interface{})           { logRaw(LevelExit, format, args...) }
func Info(format string, args ...interface{})           { logRaw(LevelInfo, format, args...) }
func Warning(format string, args ...interface{})        { logRaw(LevelWarning, format, args...) }
func Error(format string, args ...interface{})          { logRaw(LevelError, format, args...) }

