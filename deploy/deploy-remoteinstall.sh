#!/bin/bash
#deploy-installfiles.  Part of the deploy.go deploy package.

set -eu -o pipefail

returnValue=0
hostname=`hostname`

chmod -R item.attributes.permissions target
chown -R item.attributes.owner:item.attributes.group target

cp -R deploy item.targetPath

if (( $returnValue == 0 )); then
	echo "Host $hostname has all required owners and groups.  Ready for transfer."
	rm -Rf ~/deploy/
	mkdir -p ~/deploy/
else
	echo "Cannot deploy to $hostname: The host is missing some users and groups." >&2
	fi
exit $returnValue
