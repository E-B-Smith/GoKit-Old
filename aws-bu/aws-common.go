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

var kSchemaVersion string = "1.00.001"
var database *sql.DB
var psqlpath string

func connectDatabase() AWSResultCode {

	//	Find psql -- 
	var error error
	psqlpath, error = exec.LookPath("psql")
	if error != nil {
		log(AWSLogError, "Can't find Postgres 'psql': %v.", error);
		return AWSResultError;
		}
	log(AWSLogDebug, "psqlpath: %v.", psqlpath)

	//	Start the database -- 

	//	Make a connection --
	database, error = sql.Open("postgres", "user=Edward dbname=Edward sslmode=disable")
	if error != nil {
		database = nil
		log(AWSLogError, "Error: Can't open database connection: %v.", error);
		return AWSResultError
		}

	//	Make sure a compatible schema is installed -- 
	rows, error := database.Query("select version from AWSParameterTable;")
	if error != nil {
		log(AWSLogError, "Error: Can't read database version: %v.", error);
		disconnectDatabase()
		return AWSResultNotInstalled
		}
	var version string
	rows.Next()
	rows.Scan(&version)
	if version != kSchemaVersion {
		log(AWSLogError, "Error: Uncompatible database version %v.  Expected %v.", version, kSchemaVersion);
		disconnectDatabase()
		return AWSResultNotInstalled
		}	

	return AWSResultSuccess
	}


func disconnectDatabase() {
	if database != nil {
		database.Close()
		database = nil
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

	if database == nil || loggingError {
		return 
		}

	var sqlCommand string =
		fmt.Sprintf(
			`insert into AWSLogTable
					(time, processname, level, pid, message)
			values	(to_timestamp(%d), '%s', '%s'::AWSLogLevel, %d, '%s');`,
			time.Now().UTC().Unix(), command, logLevel, os.Getpid(), strings.Replace(message, "'", "''", -1))

	_, error := database.Exec(sqlCommand);
	if error != nil {
		loggingError = true
		log(AWSLogError, "Error while logging: %v.", error)
		}
	}

