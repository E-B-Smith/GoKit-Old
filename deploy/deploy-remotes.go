//  deploy  -  Deploy utility.  Deploy a set of files across a set of servers.
//
//  E.B.Smith  -  November, 2014


package main


func SSHCommand(remoteHost string, command string) (result string, error error) {
	//	ssh -o StrictHostKeyChecking=no
	
	return "", nil
	}

func CheckRemoteHostsForOwnersAndGroups(group DeployGroup) error {
	//	Load script.
	//	-	Check users in script.  % if id Edss >& /dev/null; then echo "true"; else echo "false"; fi
	// 	-	Check groups in script. % if grep "^wheel:" /etc/group > /dev/null; then echo "true"; else echo "false"; fi
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
		log(DULogDebug, "Host: %s Owners: %v Groups: %v.", host, owners, groups)
		commandString := "echo Hello!"
		SSHCommand(host, commandString)
		}

	return nil
	}


func CopyFileToRemoteHosts(group DeployGroup) error {
	//	rsync files to remote.
	//	rsync -vutzr -e 'ssh -o StrictHostKeyChecking=no' --chmod u=rwx,go-rwx --delete --force --exclude '.*' --progress --partial

	for i := 0; i < len(group.deployItems); i++ {

		item := &group.deployItems[i]
		source := item.sourcePath
		target := "~/deploy/" + item.targetPath

		for j := 0; j < len(group.deployHosts); j++ {
			remoteTarget := "deplybot@"+group.deployHosts[j]+":"+target
			log(DULogDebug, "rsync %s as %s", source, remoteTarget)
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

		for j := 0; j < len(group.deployItems); j++ {

			item := &group.deployItems[i]
			target := "~/deploy/"+item.targetPath;

			log(DULogDebug, "\tchmod -R %d %s", item.attributes.permissions, target)
			log(DULogDebug, "\tchown -R %s:%s %s", item.attributes.owner, item.attributes.group, target)
			}
		}

	for i := 0; i < len(group.deployHosts); i++ {
		host := &group.deployHosts[i]
		log(DULogDebug, "%s", *host)

		for j := 0; j < len(group.deployItems); j++ {

			item := &group.deployItems[j]
			deployPath := "~/deploy/"+item.targetPath;

			log(DULogDebug, "\tcp -R %s %s", deployPath, item.targetPath)
			}
		}

	return nil
	}

