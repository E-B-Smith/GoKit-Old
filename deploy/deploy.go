//  deploy  -  Deploy utility.  Deploy a set of files across a set of servers.
//
//  E.B.Smith  -  November, 2014

//  deploy -nfrvi [host:]/files/deplydata  relcy-r1:hduser:hduser:644:~/data
//  deploy relcy-r1:hduser:hduser:644:~/data2/file1
//
//   -n   Dry-run
//   -f   Force.
//   -r   Reverse deploy direction.
//   -v   Verbose.
//   -i   Input manifest file.


package main


import (
	"os"
	"fmt"
	"flag"
	"bufio"
	"strings"
	)


var (
	flagDryRun bool = false
	flagForceRun bool = false
	flagReverseDeploy bool = false
	flagVerbose bool = false
	flagInputFileName string
	flagInputFile *os.File = nil
	)


func main() {
	if len(os.Args) <= 1 {
		fmt.Fprintf(os.Stderr, "Error: A deploy command is expected.  Try 'deploy --help' for help.\n")
		os.Exit(1)
		}

	log(DULogDebug, "Command line: %s.", strings.Trim(fmt.Sprint(os.Args), "[]"))


	flag.BoolVar(&flagDryRun,  "n", false, "Dry run.  Just write what would be done.")
	flag.BoolVar(&flagForceRun, "f", false, "Force.  Run the deployment.")
	flag.BoolVar(&flagReverseDeploy, "r", false, "Reverse.  Reverse the direction of the deployment.")
	flag.BoolVar(&flagVerbose, "v", false, "Verbose.  Verbose output.")
	flag.StringVar(&flagInputFileName, "i", "", "Input file.  The file from which to read the deployment manifest.")
	flag.Parse()


	var error error = nil
	if flagInputFileName == "" {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {		
			flagInputFile = os.Stdin 
			}
	} else  {
		flagInputFile, error = os.Open(flagInputFileName)
		if error != nil {
            log(DULogError, "Error: Can't open file '%s' for reading: %s.", flagInputFileName, error)
            os.Exit(1)
        	}
        defer flagInputFile.Close()
		}

	log(DULogStart, "Start %s.", strings.Trim(fmt.Sprint(os.Args), "[]"))

    scanner := bufio.NewScanner(flagInputFile)
    scanner.Split(ScanManifest)
    for scanner.Scan() {
        str := scanner.Text()
        log(DULogDebug, "%s", str)
    	}
	}

