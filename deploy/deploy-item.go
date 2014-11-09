//  deploy  -  Deploy utility.  Deploy a set of files across a set of servers.
//
//  E.B.Smith  -  November, 2014


package main


type DeployAttbutes struct {
	owner string
	group string
	permissions int
	}

	
type DeployItem struct {
	attributes DeployAttbutes
	sourcePath string
	targetPath string
	}

