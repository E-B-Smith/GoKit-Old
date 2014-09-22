package main

import (
	"fmt"
	"os"
	"os/exec"	
	"strings"
	"time"
	_ "github.com/lib/pq"
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
var globalPSQLPath string
var globalAWSBackupBucket string = ""
var globalAWSAccessKeyID  string = "AKIAIUDYX3CQOEGT4OXQ"
var globalAWSAccessSecret string = "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"

type AWSParameters struct {
	schemaVersion string
	}


func connectDatabase() AWSResultCode {

	//	Find psql -- 
	var error error
	globalPSQLPath, error = exec.LookPath("psql")
	if error != nil {
		log(AWSLogError, "Can't find Postgres 'psql': %v.", error);
		return AWSResultError;
		}
	log(AWSLogDebug, "psqlpath: %v.", globalPSQLPath)

	//	Start the database -- 

	//	Make a connection --
	globalDatabase, error = sql.Open("postgres", "user=Edward dbname=Edward sslmode=disable")
	if error != nil {
		globalDatabase = nil
		log(AWSLogError, "Error: Can't open database connection: %v.", error);
		return AWSResultError
		}

	//	Make sure a compatible schema is installed -- 
	rows, error := globalDatabase.Query("select version from AWSParameterTable;")
	if error != nil {
		log(AWSLogError, "Error: Can't read database schema version: %v.", error);
		disconnectDatabase()
		return AWSResultNotInstalled
		}
	var version string
	rows.Next()
	rows.Scan(&version)
	if version != kSchemaVersion {
		log(AWSLogError, "Error: Uncompatible database schema version '%v'.  Expected '%v'.", version, kSchemaVersion);
		disconnectDatabase()
		return AWSResultNotInstalled
		}	

	return AWSResultSuccess
	}


func disconnectDatabase() {
	if  globalDatabase != nil {
		globalDatabase.Close()
		globalDatabase = nil
		}
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

var loggingError bool = false

func log(logLevel AWSLogLevel, format string, args ...interface{}) {

	var message = fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "%13s: %s\n", logLevel, message)

	if  globalDatabase == nil || loggingError {
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
		loggingError = true
		log(AWSLogError, "Error while logging: %v.", error)
		}
	}

