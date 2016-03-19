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
    //"unicode"
    "os/user"
    "path/filepath"
)


type LogLevelType int
const (
    LogLevelInvalid LogLevelType = iota
    LogLevelAll
    LogLevelDebug
    LogLevelStart
    LogLevelExit
    LogLevelInfo
    LogLevelWarning
    LogLevelError
)
var levelNames = []string {
    "LogLevelInvalid",
    "LogLevelAll",
    "LogLevelDebug",
    "LogLevelStart",
    "LogLevelExit",
    "LogLevelInfo",
    "LogLevelWarning",
    "LogLevelError",
}
var (
    LogLevel        LogLevelType    = LogLevelInfo
    logWriter       io.WriteCloser  = os.Stderr
    LogTeeStderr    bool            = false
    logFilename     string          = ""
    logRotationTime time.Time
    LogRotationInterval time.Duration = time.Hour * 24.0
)


func LogLevelFromString(s string) LogLevelType {
    for index := range levelNames {
        if s == levelNames[index] {
            return LogLevelType(index)
        }
    }
    return LogLevelInvalid
}


func StringFromLogLevel(level LogLevelType) string {
    if level < LogLevelInvalid || level > LogLevelError {
        return levelNames[LogLevelInvalid]
    } else {
        return levelNames[level]
    }
}


func homePath() string {
    homepath := ""
    u, error := user.Current()
    if error == nil {
        homepath = u.HomeDir
    } else {
        homepath = os.Getenv("HOME")
    }
    return homepath
}


func absolutePath(filename string) string {
    filename = strings.TrimSpace(filename)
    if  filepath.HasPrefix(filename, "~") {
        filename = strings.TrimPrefix(filename, "~")
        filename = path.Join(homePath(), filename)
    }
    if ! path.IsAbs(filename) {
        s, _ := os.Getwd()
        filename = path.Join(s, filename)
    }
    filename = path.Clean(filename)
    return filename
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
    logRotationTime = time.Unix(math.MaxInt64 - 1000, 0)  //  Distant future

    defer func() {
        if reason := recover(); reason != nil {
            logFilename = ""
            logWriter = os.Stderr
            Errorf("%s", reason)
        }
        name := logFilename
        if name == "" { name = "Stderr" }
        fmt.Fprintf(os.Stderr, "Log file is '%s'.\n", name)
    }()

    //fmt.Fprintf(os.Stderr, "Log: %s.\n", logFilename)

    logFilename = strings.TrimSpace(logFilename)
    if logFilename == "" {
        logWriter = os.Stderr
        return
    }

    logFilename = absolutePath(logFilename)
    if len(logFilename) <= 0 {
        logWriter = os.Stderr
        return
    }

    //fmt.Fprintf(os.Stderr, "Log: %s.\n", logFilename)

    var error error
    pathname := filepath.Dir(logFilename)
    if len(pathname) > 0 {
        if error = os.MkdirAll(pathname, 0700); error != nil {
            logWriter = os.Stderr
            panic(fmt.Sprintf("Can't create directory for log file '%s': %v.", logFilename, error))
        }
    }

    //fmt.Fprintf(os.Stderr, "Log: %s\n.", pathname)

    var flags int = syscall.O_APPEND | syscall.O_CREAT | syscall.O_WRONLY
    var mode os.FileMode = os.ModeAppend | 0700

    //fmt.Fprintf(os.Stderr, "Log: %s\n.", logFilename)

    logWriter, error = os.OpenFile(logFilename, flags, mode)
    if error != nil {
        logWriter = os.Stderr
        panic(fmt.Sprintf("Can't open log file '%s' for writing: %v.", logFilename, error))
    }

    if LogRotationInterval.Seconds() > 0 {
        var nextTime int64 = (int64(time.Now().Unix()) / int64(LogRotationInterval.Seconds())) + 1
        nextTime *= int64(LogRotationInterval.Seconds())
        logRotationTime = time.Unix(nextTime, 0)
    }
}


func rotateLogFile() {
    if len(logFilename) <= 0 { return }

    defer func() {
        if reason := recover(); reason != nil {
            logFilename = ""
            logWriter = os.Stderr
            Errorf("%s", reason)
        }
        name := logFilename
        if name == "" { name = "Stderr" }
        fmt.Fprintf(os.Stderr, "Log file is '%s'.\n", name)
    }()

/*
    replacePunct := func(r rune) rune {
        if unicode.IsLetter(r) || unicode.IsDigit(r) {
            return r
        } else {
            return '-'
        }
    }
*/

    //  Create a new file for the log --

    baseName := filepath.Base(logFilename)
    ext := filepath.Ext(baseName)
    if len(ext) != 0 {
        baseName = strings.TrimSuffix(baseName, ext)
    }
    //timeString := strings.Map(replacePunct, logRotationTime.Format(time.RFC3339)) + ext
    timeString := logRotationTime.Format(time.RFC3339) + ext
    newBase := fmt.Sprintf("%s-%s", baseName, timeString)
    newPath := filepath.Join(filepath.Dir(logFilename), newBase)
    closeLogFile()
    error := os.Rename(logFilename, newPath)
    if error != nil { panic(error) }
    openLogFile()
    Infof("Log rotated to '%s'.", newPath)
    Infof("Log continues in '%s'.", logFilename)

    //  Delete the oldest --

    globPath := filepath.Join(filepath.Dir(logFilename), baseName+"-*")
    logfiles, error := filepath.Glob(globPath)
    //fmt.Fprintf(os.Stderr, "files: %+v.\n", logfiles)
    if error != nil {
        LogError(error)
        return
    }

    //  Keep the newest 7 --

    sortedLogfiles := sort.StringSlice(logfiles)
    sortedLogfiles.Sort()
    for i := 0; i < len(sortedLogfiles) - 7; i++ {
        Infof("Removing old log '%s'.", sortedLogfiles[i])
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
    logRaw(LogLevelError, "Error '%v':\n", error)
    logRaw(LogLevelError, "Stack of %d bytes: %s\n", count, trace)
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
    if LogTeeStderr {
        fmt.Fprintf(os.Stderr, "%26s:%-4d %s: %s\n", filename[:i], linenumber, " Info", message)
    }
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
    if logLevel < LogLevelDebug || logLevel > LogLevelError { logLevel = LogLevelError }

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
    if LogTeeStderr {
        fmt.Fprintf(os.Stderr, "%26s:%-4d %s: %s\n", filename[:i], linenumber, LevelNames[logLevel], message)
    }
}


func Debugf(format string, args ...interface{})          { logRaw(LogLevelDebug, format, args...) }
func Startf(format string, args ...interface{})          { logRaw(LogLevelStart, format, args...) }
func Exitf(format string, args ...interface{})           { logRaw(LogLevelExit, format, args...) }
func Infof(format string, args ...interface{})           { logRaw(LogLevelInfo, format, args...) }
func Warningf(format string, args ...interface{})        { logRaw(LogLevelWarning, format, args...) }
func Errorf(format string, args ...interface{})          { logRaw(LogLevelError, format, args...) }
func LogError(error error)                               { logRaw(LogLevelError, "%v.", error) }

