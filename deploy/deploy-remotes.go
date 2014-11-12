//  deploy  -  Deploy utility.  Deploy a set of files across a set of servers.
//
//  E.B.Smith  -  November, 2014


package main


import (
	"io"
	"fmt"
	"sync"
	"bufio"
	"os/exec"
	)


func SSHCommand(remoteHost string, commandString string) (result string, error error) {
	//
	//	Run a string as an SSH script on a remote host -- 
	//

	if  globalSSHPath == "" {
		globalSSHPath, error = exec.LookPath("ssh")
		if error != nil {
			log(DULogError, "Can't find ssh: %v.", error)
			return "", error
			}
		}

	sshOptions := [] string {
		"-T",
		"-o", "StrictHostKeyChecking=no",
		globalDeployUser+"@"+remoteHost,
		}
	command := exec.Command(globalSSHPath, sshOptions...)
//	command.Env = append(command.Env, "PGOPTIONS=-c client_min_messages=WARNING")
	commandpipe, error := command.StdinPipe()
	if error != nil {
		log(DULogError, "Can't open pipe: %v", error)
		return "", error
		}
	
	var errorpipe *io.PipeReader;
	errorpipe, command.Stderr = io.Pipe()

	error = command.Start()
	if error != nil {
		log(DULogError, "Error running ssh: %v.", error)
		return "", error
		}

	commandpipe.Write([]byte(commandString))
	commandpipe.Close()

	var waiter sync.WaitGroup
	waiter.Add(1)
	go func() {
		scanner := bufio.NewScanner(errorpipe)
		for scanner.Scan() {
			log(DULogError, "%v.", scanner.Text())
			}
		waiter.Done()
		} ()

	error = command.Wait()
	errorpipe.Close()
	waiter.Wait()

	if error != nil {
		log(DULogError, "Script %v.", error)
		return "", error
		}

	return "", nil
	}

func CheckRemoteHostsForOwnersAndGroups(group DeployGroup) error {
	//	Load script.
	//	-	Check users in script.  
	// 	-	Check groups in script.
	//  -	Script: Remove deploy directory.
	//	-	Script: Create deploy directory.

	owners := make(map[string]bool)
	groups := make(map[string]bool)

	for i := 0; i < len(group.deployItems); i++ {
		attributes := &group.deployItems[i].attributes
		owners[attributes.owner] = true;
		groups[attributes.group] = true;
		}

	for i := 0; i < len(group.deployHosts); i++ {
		host := group.deployHosts[i]

		commandString :=
			`#!/bin/bash
			 #deploy-remotecheck.  Part of the deploy.go deploy package.

			set -eu -o pipefail

			returnValue=0
			hostname=$(hostname)
			function ckowner()
				{
				if ! id "$1" >& /dev/null; then 
					echo "Host $hostname user $1 does not exist" >&2
					returnValue=1
					fi
				}
			function ckgroup()
				{
			 	if ! grep "^$1:" /etc/group >& /dev/null; then 
					echo "Host $hostname group $1 does not exist" >&2
					returnValue=1 		
			 		fi	
				}
			`

		for owner, _ := range owners {
			commandString += fmt.Sprintf("ckowner %s\n", owner)
			}

		for group, _ := range owners {
			commandString += fmt.Sprintf("ckgroup %s\n", group)
			}

		commandString += 
			`
			if (( $returnValue == 0 )); then
				echo "Host $hostname has all required owners and groups.  Ready for transfer."
				rm -Rf ~/deploy/
				mkdir -p ~/deploy/
			else
				echo "Cannot deploy to $hostname: The host is missing some users and groups." >&2
				fi
			exit $returnValue
			`
		SSHCommand(host, commandString)
		}

	return nil
	}


func RsyncCommand(sourcePath string, targetPath string) (result string, error error) {
	//
	//	Rsync file to a remote host.
	//
	//	rsync -vtzr -e 'ssh -o StrictHostKeyChecking=no' --chmod u=rwx,go-rwx 
	//	      --exclude '.DS_Store' --partial
	
	if  globalRsyncPath == "" {
		globalRsyncPath, error = exec.LookPath("rsync")
		if error != nil {
			log(DULogError, "Can't find rsync: %v.", error)
			return "", error
			}
		}

	rsyncOptions := [] string {
		"-vtzr",
		"-e", "ssh -o StrictHostKeyChecking=no",
		"--chmod", "u=rwx,go-rwx",
		"--exclude", ".DS_Store",
		"--partial",
		sourcePath, targetPath,
		}
	command := exec.Command(globalRsyncPath, rsyncOptions...)	
	var errorpipe *io.PipeReader;
	errorpipe, command.Stderr = io.Pipe()

	error = command.Start()
	if error != nil {
		log(DULogError, "Error running rsync: %v.", error)
		return "", error
		}

	var waiter sync.WaitGroup
	waiter.Add(1)
	go func() {
		scanner := bufio.NewScanner(errorpipe)
		for scanner.Scan() {
			log(DULogError, "%v.", scanner.Text())
			}
		waiter.Done()
		} ()

	error = command.Wait()
	errorpipe.Close()
	waiter.Wait()

	if error != nil {
		log(DULogError, "Script %v.", error)
		return "", error
		}

	return "", nil
	}


func CopyFileToRemoteHosts(group DeployGroup) error {
	//	Copy each file to the remote host.

	for i := 0; i < len(group.deployItems); i++ {

		item := &group.deployItems[i]
		source := item.sourcePath
		target := "~/deploy/" + item.targetPath

		for j := 0; j < len(group.deployHosts); j++ {
			remoteTarget := globalDeployUser+"@"+group.deployHosts[j]+":"+target
			RsyncCommand(source, remoteTarget)
			}
		}

	return nil
	}


func InstallFilesOnRemoteHosts(group DeployGroup) error {
	//	Load script.
	//	Add deploy items to script.
	//	-	SetPermissions + owners of every item.
	//	-	Move items into place.


	for i := 0; i < len(group.deployHosts); i++ {
		host := &group.deployHosts[i]
		log(DULogDebug, "%s", *host)

		commandString := 
			`#!/bin/bash
			 #deploy-installfiles.  Part of the deploy.go deploy package.

			set -eu -o pipefail

			returnValue=0
			hostname=$(hostname)
			`

		for j := 0; j < len(group.deployItems); j++ {

			item := &group.deployItems[i]
			target := "~/deploy/"+item.targetPath;

			commandString += fmt.Sprintf("chmod -R %d %s\n", item.attributes.permissions, target)
			commandString += fmt.Sprintf("chown -R %s:%s %s\n", item.attributes.owner, item.attributes.group, target)
			}

		for j := 0; j < len(group.deployItems); j++ {

			item := &group.deployItems[j]
			deployPath := "~/deploy/"+item.targetPath;

			commandString += fmt.Sprintf("cp -R %s %s\n", deployPath, item.targetPath)
			}

		commandString +=
			`
			if (( $returnValue == 0 )); then
				echo "Host $hostname has all required owners and groups.  Ready for transfer."
				rm -Rf ~/deploy/
				mkdir -p ~/deploy/
			else
				echo "Cannot deploy to $hostname: The host is missing some users and groups." >&2
				fi
			exit $returnValue
			`

		SSHCommand(*host, commandString)
		}

	return nil
	}

