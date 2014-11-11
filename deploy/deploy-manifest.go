//  deploy  -  Deploy utility.  Deploy a set of files across a set of servers.
//
//  E.B.Smith  -  November, 2014


package main


import (
	"io"
	"os"
	"fmt"
	"path"
	"errors"
	)


type DeployAttributes struct {
	owner string
	group string
	permissions int
	}

	
type DeployItem struct {
	attributes DeployAttributes
	sourcePath string
	targetPath string
	}


type DeployGroup struct {
	deployItems [] DeployItem
	deployHosts [] string
	}


type DeployManifest struct {
	deployGroups [] DeployGroup
	}


func ParseError(scanner *ZScanner, message string) error {
	basename := path.Base(scanner.FileName())
	message = 
		fmt.Sprintf("%s:%d Scanned '%s'. %s",
			basename, scanner.LineNumber(), scanner.Token(), message)
	return errors.New(message)
	}


func ParseHosts(scanner *ZScanner) (hosts []string, error error) {
	host, error := scanner.ScanIdentifier()
	log(DULogDebug, "%s %v", host, error)
	for error == nil && host != ";" {
		hosts = append(hosts, host)
		host, error = scanner.ScanNext()
		log(DULogDebug, "%s %v", host, error)
		if host == "," { 
			host, error = scanner.ScanNext()
			}
		}

	return hosts, error
	}


func ParseDeployItems(scanner *ZScanner) (newItem []DeployItem, error error) {

	var attributes DeployAttributes
	group := new(DeployGroup)

	token, error := scanner.ScanNext()
	if token != "{" {
		return nil, ParseError(scanner, "A deploy-group starts with a bracket '{'.")
		}

	token, error = scanner.ScanNext()
	for !scanner.IsAtEnd() && token != "}" {

		log(DULogDebug, "Scanned '%s'.", scanner.token)

		if token == "owner" {
			if attributes.owner != "" {
				return nil, ParseError(scanner, "An owner was already specified")
				}
			token, _ = scanner.ScanIdentifier()
			if token == "" {
				return nil, ParseError(scanner, "A user name is expected")
				}
			attributes.owner = token
			token, _ = scanner.ScanNext()
			if token != ";" {
				return nil, ParseError(scanner, "A semi-colon is expected")
				}
			token, error = scanner.ScanNext()
			continue
			}


		if token == "group" {
			if attributes.group != "" {
				return nil, ParseError(scanner, "A group was already specified")
				}
			token, _ = scanner.ScanIdentifier()
			if token == "" {
				return nil, ParseError(scanner, "A group name is expected")
				}
			attributes.group = token
			token, _ = scanner.ScanNext()
			if token != ";" {
				return nil, ParseError(scanner, "A semi-colon is expected")
				}
			token, error = scanner.ScanNext()
			continue
			}


		if token == "permissions" {
			if attributes.permissions != 0 {
				return nil, ParseError(scanner, "Permissions were already specified")
				}
			mode, error := scanner.ScanInteger()
			if error != nil {
				return nil, ParseError(scanner, "Permissions are expected")
				}
			if mode == 0 {
				return nil, ParseError(scanner, "Permissions of mode '000' aren't permitted")
				}				
			attributes.permissions = mode
			token, _ = scanner.ScanNext()
			if token != ";" {
				return nil, ParseError(scanner, "A semi-colon is expected")
				}
			token, error = scanner.ScanNext()
			continue
			}


		if token == "deploy" {
			var deployItem DeployItem
			deployItem.sourcePath, error = scanner.ScanString()
			if error != nil {
				return nil, ParseError(scanner, "Expected a source file name")
				}

			scanner.ScanNext()
			if scanner.error != nil {
				return nil, ParseError(scanner, "The end of a deploy group was expected")
				}

			if scanner.token == "as" {
				deployItem.targetPath, error = scanner.ScanString()
				if error != nil {
					return nil, ParseError(scanner, "Expected a target file name")
					}
				scanner.ScanNext()
				}

			if scanner.token != ";" {
				return nil, ParseError(scanner, "A semi-colon is expected at the end of a deploy statement")
				}

			group.deployItems = append(group.deployItems, deployItem)

			token, error = scanner.ScanNext()
			continue
			}


		if token == "deploy-group" {
			deployItems, error := ParseDeployItems(scanner)
			if error != nil {
				return nil, error
				}
			group.deployItems = append(group.deployItems, deployItems...)
			token, error = scanner.ScanNext()
			continue
			}			

		return nil, ParseError(scanner, "Unrecognized keyword")
		}

	//	Apply the attributes -- 

	for i:=0; i < len(group.deployItems); i++ {
		item := &group.deployItems[i]
		if attributes.owner != "" {
			if item.attributes.owner == "" { item.attributes.owner = attributes.owner }
			}
		if attributes.group != "" {
			if item.attributes.group == "" { item.attributes.group = attributes.group }
			}
		if attributes.permissions != 0 {
			if item.attributes.permissions == 0 { item.attributes.permissions = attributes.permissions }
			}
		log(DULogDebug, "Item: %v", item)
		}

	return group.deployItems, nil 
	}


func ParseManifest(inputFile *os.File) (manifest *DeployManifest, error error) {

	manifest = new(DeployManifest)
	var currentGroup *DeployGroup
	currentGroup = nil

	scanner := NewZScanner(inputFile)
	for !scanner.IsAtEnd() {
		var indentifier string

		indentifier, error = scanner.ScanIdentifier()
		log(DULogDebug, "Scanned '%s'.", scanner.token)

		if error == io.EOF {
			return manifest, nil
			}
		if error != nil {
		    return nil, error
			}
		
		if indentifier == "deploy-group" {
		    deployItems, error := ParseDeployItems(scanner)
		    if error != nil { return nil, error }
		    if (currentGroup == nil) {
		    	currentGroup = new(DeployGroup)
		    	}
		    currentGroup.deployItems = append(currentGroup.deployItems, deployItems...)
		    continue
		    }


		if indentifier == "hosts" {
		    currentGroup.deployHosts, error = ParseHosts(scanner)
		    if error != nil { return nil, error }
		    manifest.deployGroups = append(manifest.deployGroups, *currentGroup)
		    currentGroup = nil
		    continue
		 	}

		return nil, ParseError(scanner, "'deploy-group' or 'hosts' expected.")
		}

	return manifest, nil
	}	

