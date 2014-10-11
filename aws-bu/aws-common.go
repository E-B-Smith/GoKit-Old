package main

import (
	"fmt"
	"os"
	"time"
	"strings"
	"database/sql"	
	)


type AWSResultCode int 
const (
	AWSResultSuccess = 0
	AWSResultWarning = 1
	AWSResultError   = 2
	AWSResultNotInstalled = 3
	)


const kSchemaVersion = "1.00.001"

var globalDatabase *sql.DB = nil
var globalPGCTLPath       string
var globalPSQLPath        string
var globalPSQLDataPath    string = "~/Library/Application Support/Postgres/var-9.3"
var globalAWSBackupBucket string = "rookery.backup"
var globalAWSAccessKeyID  string = "AKIAIUDYX3CQOEGT4OXQ"
var globalAWSAccessSecret string = "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"
var globalAWSRegion 	  string = "us-west-2"


type AWSParameters struct {
	schemaVersion string
	}


type AWSLogLevel string
const (
	AWSLogDebug   = "AWSLogDebug"
	AWSLogInfo	  = "AWSLogInfo"
	AWSLogStart   = "AWSLogStart"
	AWSLogExit    = "AWSLogExit"
	AWSLogWarning = "AWSLogWarning"
	AWSLogError   = "AWSLogError"
	)


var globalLoggingError bool = false

func log(logLevel AWSLogLevel, format string, args ...interface{}) {

	var message = fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "%13s: %s\n", logLevel, message)

	if  globalDatabase == nil || globalLoggingError {
		return 
		}

	var sqlCommand string =
		fmt.Sprintf(
			`insert into AWSLogTable
					(time, processname, level, pid, message)
			values	(to_timestamp(%d), '%s', '%s'::AWSLogLevel, %d, '%s');`,
			time.Now().UTC().Unix(), command, logLevel, os.Getpid(), strings.Replace(message, "'", "''", -1))

	_, error := globalDatabase.Exec(sqlCommand);
	if error != nil {
		globalLoggingError = true
		log(AWSLogError, "Error while logging: %v.", error)
		}
	}

