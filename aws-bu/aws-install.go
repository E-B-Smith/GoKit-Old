package main

import (
//	"os"
//	"fmt"
//	"log"
	"os/exec"
//	"github.com/jteeuwen/go-bindata"
	)
	

func install() int {
	psqlPath, error := exec.LookPath("psql")
	if error != nil {
		log(AWSLogError, "Can't find and execute 'psql'.");
		return kResultError;
		}
	
	var sqlInstallScript [] byte;
	sqlInstallScript, error = Asset("install.sql")
	if len(sqlInstallScript) == 0 {
		log(AWSLogError, "Can't load asset 'install.sql'.");
		return kResultError;
		}
	
	psqlOptions := [] string {
		"-h", "localhost",
		"-X", "-q",
		"-v", "ON_ERROR_STOP=1",
		"--pset", "pager=off",
		}
	command := exec.Command(psqlPath, psqlOptions...)
	command.Run()
	
	return kResultSuccess;
	}