//  log  -  Simple logging.
//
//  E.B.Smith  -  November, 2014


package log


import (
    "io"
    "os"
    "fmt"
    "path"
    "syscall"
    "runtime"
    "path/filepath"
)


type LogLevelType int
const (
    LevelInvalid LogLevelType = iota
    LevelAll
    LevelDebug
    LevelStart
    LevelExit
    LevelInfo
    LevelWarning
    LevelError
)
var levelNames = []string {
    "LevelInvalid",
    "LevelAll",
    "LevelDebug",
    "LevelStart",
    "LevelExit",
    "LevelInfo",
    "LevelWarning",
    "LevelError",
}


var LogLevel    LogLevelType    = LevelWarning
var logWriter   io.WriteCloser  = os.Stderr


func LogLevelFromString(s string) LogLevelType {
    for index := range levelNames {
        if s == levelNames[index] {
            return LogLevelType(index)
        }
    }
    return LevelInvalid
}


func StringFromLogLevel(level LogLevelType) string {
    if level < LevelInvalid || level > LevelError {
        return levelNames[LevelInvalid]
    } else {
        return levelNames[level]
    }
}


func SetFilename(filename string) {
    _, hasClose := logWriter.(interface {Close()})
    if  hasClose &&
        logWriter != os.Stderr &&
        logWriter != os.Stdout {
        logWriter.Close()
    }
    if len(filename) <= 0 {
        logWriter = os.Stderr
        return
    }
    var error error
    var flags int = syscall.O_APPEND | syscall.O_CREAT | syscall.O_WRONLY
    var mode os.FileMode = os.ModeAppend | 0700
    pathname := filepath.Dir(filename)
    if len(pathname) > 0 {
        if error = os.MkdirAll(pathname, 0700); error != nil {
            logWriter = os.Stderr
            Errorf("Error: Can't create directory for log file '%s': %v.", filename, error)
        }
    }
    logWriter, error = os.OpenFile(filename, flags, mode)
    if error != nil {
        logWriter = os.Stderr
        Errorf("Error: Can't open log file '%s' for writing: %v.", filename, error)
    }
}


func LogStackWithError(error interface{}) {
    trace := make([]byte, 1024)
    count := runtime.Stack(trace, true)
    logRaw(LevelError, "Error '%v':\n", error)
    logRaw(LevelError, "Stack of %d bytes: %s\n", count, trace)
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
    fmt.Fprintf(logWriter, "%26s:%-4d %s: %s\n", filename[:i], linenumber, " Info", message)
}


func FlushMessages() {
    logWriter.Close()
}


func logRaw(logLevel LogLevelType, format string, args ...interface{}) {

    LevelNames := []string {
        "Inval",
        "  All",
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
    fmt.Fprintf(logWriter, "%26s:%-4d %s: %s\n", filename[:i], linenumber, LevelNames[logLevel], message)
}


func Debugf(format string, args ...interface{})          { logRaw(LevelDebug, format, args...) }
func Startf(format string, args ...interface{})          { logRaw(LevelStart, format, args...) }
func Exitf(format string, args ...interface{})           { logRaw(LevelExit, format, args...) }
func Infof(format string, args ...interface{})           { logRaw(LevelInfo, format, args...) }
func Warningf(format string, args ...interface{})        { logRaw(LevelWarning, format, args...) }
func Errorf(format string, args ...interface{})          { logRaw(LevelError, format, args...) }
func LogError(error error)                              { logRaw(LevelError, "%v.", error) }

