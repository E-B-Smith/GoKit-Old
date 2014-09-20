package main

import (
	"io"
	"sync"
	"bufio"
	"os/exec"
	)
	

func runSQLScript(scriptname string) AWSResultCode {
	
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
	command := exec.Command(psqlpath, psqlOptions...)
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


func install() AWSResultCode {
	//	Check if already installed --
	//	If installed, ask for reinstall.

	//	Install -- 
	result := runSQLScript("aws-resources/install.sql")
	if result != AWSResultSuccess {
		return result;
		}

	connectDatabase()
	log(AWSLogStart, "Start install.");

	return AWSResultSuccess;
	}


func uninstall() AWSResultCode {	
	result := runSQLScript("aws-resources/uninstall.sql")
	disconnectDatabase()
	log(AWSLogInfo, "aws-bu uninstalled.")
	return result
	}

