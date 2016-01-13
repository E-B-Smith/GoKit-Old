//  Log  -  Simple logging.
//
//  E.B.Smith  -  November, 2014


package Log


import (
    "io"
    "os"
    "fmt"
    "math"
    "path"
    "sort"
    "time"
    "runtime"
    "syscall"
    "strings"
    "unicode"
    "path/filepath"
    "violent.blue/GoKit/Util"
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
var (
    LogLevel        LogLevelType    = LevelInfo
    logWriter       io.WriteCloser  = os.Stderr
    logFilename     string          = ""
    logRotationTime time.Time
    logRotationInterval time.Duration = time.Hour * 24.0
)


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


func closeLogFile() {
    _, hasClose := logWriter.(interface {Close()})
    if  hasClose &&
        logWriter != os.Stderr &&
        logWriter != os.Stdout {
        logWriter.Close()
    }
}


func openLogFile() {
    logRotationTime = time.Unix(math.MaxInt64 - 10000, 0)  //  Distant future

    logFilename = ServerUtil.AbsolutePath(logFilename)
    if len(logFilename) <= 0 {
        logWriter = os.Stderr
        return
    }

    var error error
    pathname := filepath.Dir(logFilename)
    if len(pathname) > 0 {
        if error = os.MkdirAll(pathname, 0700); error != nil {
            logWriter = os.Stderr
            Errorf("Error: Can't create directory for log file '%s': %v.", logFilename, error)
        }
    }

    var flags int = syscall.O_APPEND | syscall.O_CREAT | syscall.O_WRONLY
    var mode os.FileMode = os.ModeAppend | 0700

    logWriter, error = os.OpenFile(logFilename, flags, mode)
    if error != nil {
        logWriter = os.Stderr
        Errorf("Error: Can't open log file '%s' for writing: %v.", logFilename, error)
    }

    if logRotationInterval.Seconds() > 0 {
        var nextTime int64 = (int64(time.Now().Unix()) / int64(logRotationInterval.Seconds())) + 1
        nextTime *= int64(logRotationInterval.Seconds())
        logRotationTime = time.Unix(nextTime, 0)
    }
}


func rotateLogFile() {
    if len(logFilename) <= 0 { return }

    replacePunct := func(r rune) rune {
        if unicode.IsLetter(r) || unicode.IsDigit(r) {
            return r
        } else {
            return '-'
        }
    }

    //  Create a new file for the log --

    closeLogFile()
    timeString := strings.Map(replacePunct, logRotationTime.Format(time.RFC3339))
    newPath := fmt.Sprintf("%s-%s", logFilename, timeString)
    error := os.Rename(logFilename, newPath)
    if error != nil { panic(error) }
    openLogFile()

    //  Delete the oldest --

    logfiles, error := filepath.Glob(logFilename+"-*")
    if error != nil {
        LogError(error)
        return
    }

    //  Keep the newest 7 --

    sortedLogfiles := sort.StringSlice(logfiles)
    sortedLogfiles.Sort()
    for i := 0; i < len(sortedLogfiles) - 7; i++ {
        error = os.Remove(sortedLogfiles[i])
        if error != nil {
            Errorf("Can't remove log file '%s': %v.", sortedLogfiles[i], error)
        }
    }
}


func SetFilename(filename string) {
    if filename == logFilename {
        return
    }
    closeLogFile()
    logFilename = filename
    openLogFile()
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
    closeLogFile()
    openLogFile()
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

    itemTime := time.Now()
    if  itemTime.After(logRotationTime)  {
        rotateLogFile()
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
func LogError(error error)                               { logRaw(LevelError, "%v.", error) }

