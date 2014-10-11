package main

import (
	"io"
	"sync"
	"bufio"
	"os/exec"
	_ "github.com/lib/pq"
	"database/sql"	
	)


func connectDatabase() AWSResultCode {

	//
	//	Start the database -- 
	//

	//	Find postgres --
	var error error
	globalPGCTLPath, error = exec.LookPath("pg_ctl")
	if error != nil {
		log(AWSLogError, "Can't find Postgres 'pg_ctl': %v.", error);
		return AWSResultError;
		}
	log(AWSLogDebug, " Found pg_ctl: %v.", globalPGCTLPath)
	log(AWSLogDebug, "Database Data: %v.", globalPSQLDataPath)

	command := exec.Command(globalPGCTLPath, "status", "-D", globalPSQLDataPath)
	error = command.Run()
	if command.ProcessState.Sys() == 3 {
		log(AWSLogDebug, "Starting Postgres")

		command = exec.Command(globalPGCTLPath, "start", "-w", "-s", "-D", globalPSQLDataPath)
		error = command.Run()
		if error != nil {
			log(AWSLogError, "Error starting Postgress: %v.", error)
			return AWSResultError
			}
	} else {
		log(AWSLogDebug, "Postgres is already started.")
		}


	//
	//	Find psql command line utility and connect -- 
	//


	//	Find psql -- 
	globalPSQLPath, error = exec.LookPath("psql")
	if error != nil {
		log(AWSLogError, "Can't find Postgres 'psql': %v.", error);
		return AWSResultError;
		}
	log(AWSLogDebug, "psqlpath: %v.", globalPSQLPath)

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
	

func runSQLScript(scriptname string) AWSResultCode {
	
	//
	//	Run an SQL script that is stored as a resource -- 
	//

	var error error
	var sqlInstallScript [] byte
	sqlInstallScript, error = Asset(scriptname)
	if len(sqlInstallScript) == 0 {
		log(AWSLogError, "Can't load asset: %v.", error);
		return AWSResultError;
		}

	psqlOptions := [] string {
		"-h", "localhost",
		"-X", "-q",
		"-v", "ON_ERROR_STOP=1",
		"--pset", "pager=off",
		}
	command := exec.Command(globalPSQLPath, psqlOptions...)
	command.Env = append(command.Env, "PGOPTIONS=-c client_min_messages=WARNING")
	commandpipe, error := command.StdinPipe()
	if error != nil {
		log(AWSLogError, "Can't open pipe: %v", error)
		return AWSResultError
		}
	
	var errorpipe *io.PipeReader;
	errorpipe, command.Stderr = io.Pipe()

	error = command.Start()
	if error != nil {
		log(AWSLogError, "Error running psql: %v.", error)
		return AWSResultError
		}

	commandpipe.Write(sqlInstallScript)
	commandpipe.Close()

	var waiter sync.WaitGroup
	waiter.Add(1)
	go func() {
		scanner := bufio.NewScanner(errorpipe)
		for scanner.Scan() {
			log(AWSLogError, "%v.", scanner.Text())
			}
		waiter.Done()
		} ()

	error = command.Wait()
	errorpipe.Close()
	waiter.Wait()

	if error != nil {
		log(AWSLogError, "Script %v.", error)
		return AWSResultError
		}

	return AWSResultSuccess
	}


