//  log  -  Simple logging.
//
//  E.B.Smith  -  November, 2014


package log


import (
	"fmt"
	"os"
	"runtime"
	"path"
	)


type LogLevel int
const (
	LevelDebug LogLevel = iota
	LevelStart  
	LevelExit   
	LevelInfo   
	LevelWarning
	LevelError
	)


var MinLogLevel = LevelDebug


func logRaw(logLevel LogLevel, format string, args ...interface{}) {

    LevelNames := []string { 
        "Debug",
        "Start",
        " Exit",
        " Info",
        " Warn",
        "Error", 
        }

    if logLevel < MinLogLevel { return }
	if logLevel < LevelDebug || logLevel > LevelError { logLevel = LevelError }

 	_, filename, linenumber, _ := runtime.Caller(2)
 	filename = path.Base(filename)
 	i := len(filename)
	if i > 16 {
 		i = 16
 		}

	var message = fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "%16s:%-4d %s %s\n", filename[:i], linenumber, LevelNames[logLevel], message)
	}


func Debug(format string, args ...interface{}) 			{ logRaw(LevelDebug, format, args...) }
func Start(format string, args ...interface{}) 			{ logRaw(LevelStart, format, args...) }
func Exit(format string, args ...interface{}) 			{ logRaw(LevelExit, format, args...) }
func Info(format string, args ...interface{}) 			{ logRaw(LevelInfo, format, args...) }
func Warning(format string, args ...interface{}) 		{ logRaw(LevelWarning, format, args...) }
func Error(format string, args ...interface{}) 			{ logRaw(LevelError, format, args...) }

