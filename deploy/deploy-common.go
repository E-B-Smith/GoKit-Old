//  deploy  -  deploy a set of files across a set of servers.
//
//  E.B.Smith  -  November, 2014


package main


import (
	"fmt"
	"os"
	"runtime"
	"path"
	)


type DUResultCode int 
const (
	DUResultSuccess = 0
	DUResultWarning = 1
	DUResultError   = 2
	DUResultNotInstalled = 3
	)


const DUVersion = "1.00.001"


type DULogLevel string
const (
	DULogDebug   = "Debug"
	DULogInfo	 = " Info"
	DULogStart   = "Start"
	DULogExit    = " Exit"
	DULogWarning = " Warn"
	DULogError   = "Error"
	)


var globalLoggingError bool = false


func log(logLevel DULogLevel, format string, args ...interface{}) {

 	_, filename, linenumber, _ := runtime.Caller(1)
 	filename = path.Base(filename)
 	i := len(filename)
	if i > 16 {
 		i = 16
 		}

	var message = fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "%16s:%-4d %s %s\n", filename[:i], linenumber, logLevel, message)
	}

