//  z/log  -  Simple logging.
//
//  E.B.Smith  -  November, 2014


package log


import (
	"fmt"
	"os"
	"runtime"
	"path"
	)


type MessageLevel int
const (
	LevelDebug MessageLevel = iota
	LevelStart  
	LevelExit   
	LevelInfo   
	LevelWarning
	LevelError
	)
	

func LogRaw(messageLevel MessageLevel, format string, args ...interface{}) {

    MessageLevelNames := []string { 
        "Debug",
        "Start",
        " Exit",
        " Info",
        " Warn",
        "Error", 
        }

	if messageLevel < 0 || messageLevel > LevelError { messageLevel = LevelError }

 	_, filename, linenumber, _ := runtime.Caller(1)
 	filename = path.Base(filename)
 	i := len(filename)
	if i > 16 {
 		i = 16
 		}

	var message = fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "%16s:%-4d %s %s\n", filename[:i], linenumber, MessageLevelNames[messageLevel], message)
	}


func Debug(format string, args ...interface{}) 			{ LogRaw(LevelDebug, format, args) }
func Start(format string, args ...interface{}) 			{ LogRaw(LevelStart, format, args) }
func Exit(format string, args ...interface{}) 			{ LogRaw(LevelExit, format, args) }
func Info(format string, args ...interface{}) 			{ LogRaw(LevelInfo, format, args) }
func Warn(format string, args ...interface{}) 			{ LogRaw(LevelWarning, format, args) }
func Error(format string, args ...interface{}) 			{ LogRaw(LevelError, format, args) }

