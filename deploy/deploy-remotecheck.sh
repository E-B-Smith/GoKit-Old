#!/bin/bash
#deploy-remotecheck.  Part of the deploy.go deploy package.

set -eu -o pipefail

returnValue=0
hostname=`hostname`
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

ckowner Edward
ckowner Bob
ckgroup staff
ckgroup blafs

if (( $returnValue == 0 )); then
	echo "Host $hostname has all required owners and groups.  Ready for transfer."
	rm -Rf ~/deploy/
	mkdir -p ~/deploy/
else
	echo "Cannot deploy to $hostname: The host is missing some users and groups." >&2
	fi
exit $returnValue
