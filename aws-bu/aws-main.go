package main

import (
	"os"
	"fmt"
	"flag"
	"strings"
	)
 

func printUsage() AWSResultCode {
    fmt.Fprintf(os.Stderr, "Usage: aws-bu [ help | list | log | history | status | report ] [ <options> ]\n\noptions:\n")
    flag.PrintDefaults()
    return 0
	}

var command = "aws-bu";

var (
	flagReverseSort bool = false
	flagFollowOutput bool = false
	flagOutputLimit int = 0
	)


func main() {

	var exitResult AWSResultCode = AWSResultError;
	exitResult = connectDatabase()
	if exitResult == AWSResultNotInstalled {
	} else 
	if exitResult != AWSResultSuccess {
		os.Exit(int(exitResult))
		}

	exitResult = AWSResultError;
	defer func() {
		log(AWSLogExit, "Exit status %d.", exitResult)
		disconnectDatabase()
		os.Exit(int(exitResult))
		} ()

	if len(os.Args) <= 1 {
		fmt.Fprintf(os.Stderr, "Error: aws-bu command expected.  Try 'aws-bu help' for help.\n")
		return
		}

//	log(AWSLogDebug, "Command line: %s.", os.shift.pop)// strings.Trim(fmt.Sprint(os.Args), "[]"))

	command = strings.Title(os.Args[1])
	log(AWSLogDebug, "Command: %v.", command)
	os.Args = os.Args[1:]

	flag.BoolVar(&flagReverseSort, "r", false, "Reverse output sort.")
	flag.BoolVar(&flagFollowOutput, "f", false, "Follow output.")
	flag.IntVar(&flagOutputLimit, "n", 0, "Number of lines to output. A zero value reports all lines.")
	flag.Parse()

	// fmt.Println("Reverse:", flagReverseSort, "\t|")
	// fmt.Println(" Follow:", flagFollowOutput, "\t|")
	// fmt.Println("  Limit:", flagOutputLimit, "\t|")

	log(AWSLogStart, "Start %s.", strings.Trim(fmt.Sprint(os.Args), "[]"))

	switch command {
		case "Help", "-h":
			exitResult = printUsage()
						
		case "List":
			exitResult = listBundles()
			
		case "Log":
			exitResult = printLog()
			
		case "History":
			exitResult = printHistory()
			
		case "Status":
			exitResult = printStatus()
		
		case "Report":
			exitResult = printReport()
		
		case "Install":
			exitResult = install()

		case "Uninstall":
			exitResult = uninstall()
		
		case "Refresh":
			exitResult = refreshData()

		default:
			fmt.Printf("Error: Unrecognized command '%v'.\n", command)
			exitResult = AWSResultError
		}
	
	//	The end.
	}

