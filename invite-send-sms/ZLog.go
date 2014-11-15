//  zlog  -  Logging,
//
//  E.B.Smith  -  November, 2014


package main


import (
	"fmt"
	"os"
	"runtime"
	"path"
	)


type ZLogLevel string
const (
	ZLogDebug   = "Debug"
	ZLogInfo	 = " Info"
	ZLogStart   = "Start"
	ZLogExit    = " Exit"
	ZLogWarning = " Warn"
	ZLogError   = "Error"
	)


func ZLog(logLevel ZLogLevel, format string, args ...interface{}) {
	//	Log a message -- 

 	_, filename, linenumber, _ := runtime.Caller(1)
 	filename = path.Base(filename)
 	i := len(filename)
	if i > 16 {
 		i = 16
 		}

	var message = fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "%16s:%-4d %s %s\n", filename[:i], linenumber, logLevel, message)
	}

