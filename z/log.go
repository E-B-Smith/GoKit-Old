//  deploy  -  deploy a set of files across a set of servers.
//
//  E.B.Smith  -  November, 2014


package log


import (
	"fmt"
	"os"
	"runtime"
	"path"
//	_ "github.com/go-sql-driver/mysql"
//	"violent.blue/go/log"
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


MessageLevelNames := []string { 
	 "Debug"
	,"Start"
	," Exit"
	," Info"
	," Warn"
	,"Error"
	}


func LogRaw(messageLevel MessageLevel, format string, args ...interface{}) {

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


func Debug(format string, args ...interface{}) {
	LogRaw(LevelDebug, format, args)
	}


log.Debug()
log.Info("This is my info message: %v.", error)
log.Warning()

struct log {

	}
