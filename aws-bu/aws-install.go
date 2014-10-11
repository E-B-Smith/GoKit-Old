package main

import (
	"io"
	"sync"
	"bufio"
	"os/exec"
	)
	

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


