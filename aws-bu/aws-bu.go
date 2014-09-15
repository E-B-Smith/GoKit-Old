package main

import (
	"os"
	"fmt"
	"flag"
	"time"
	"strings"
	_ "github.com/lib/pq"
	"database/sql"
)

type AWSLogLevel string
const (
	AWSLogDebug   = "AWSLogDebug"
	AWSLogInfo	  = "AWSLogInfo"
	AWSLogStart   = "AWSLogStart"
	AWSLogExit    = "AWSLogExit"
	AWSLogWarning = "AWSLogWarning"
	AWSLogError   = "AWSLogError"
)

var database *sql.DB;
var command = "aws-bu";
 
func log(logLevel AWSLogLevel, format string, args ...interface{}) {

	var message = fmt.Sprintf(format, args...);
	var terminalMessage = fmt.Sprintf("%13s: %s", logLevel, message)
	fmt.Println(terminalMessage)

	var sqlCommand string =
		fmt.Sprintf(
			`insert into AWSLogTable
					(time, processname, level, pid, message)
			values	(to_timestamp(%d), '%s', '%s'::AWSLogLevel, %d, '%s');`,
			time.Now().UTC().Unix(), command, logLevel, os.Getpid(), strings.Replace(message, "'", "''", -1));

	_, error := database.Exec(sqlCommand);
	if error != nil {
		fmt.Printf("Error while logging: %v.\n", error);
		fmt.Printf("Error while logging:\n%v.\n", sqlCommand);
		os.Exit(1);
	}
}

func printUsage() int {
    fmt.Fprintf(os.Stderr, "Usage: aws-bu [ help | list | log | history | status | report ] [ <options> ]\n\noptions:\n")
    flag.PrintDefaults()
    return 0
	}

var (
	flagReverseSort bool = false
	flagFollowOutput bool = false
	flagOutputLimit int = 0
	)

func main() {
	var error error;
	database, error = sql.Open("postgres", "user=Edward dbname=Edward sslmode=disable")
	if error != nil {
		fmt.Printf("Error: Can't open database connection: %v.\n", error);
		os.Exit(1)
	}
	
	log(AWSLogDebug, "Arguments: %v.", os.Args);

	if len(os.Args) <= 1 {
		fmt.Printf("Error: aws-bu command expected.  Try 'aws-bu help' for help.\n")
		os.Exit(1)
	}

	rawCommand := os.Args[1];
	command = strings.Title(rawCommand)
	log(AWSLogDebug, "Command: %v.", command)
	os.Args = os.Args[1:]

	var status int = 0;
	defer func() {
		log(AWSLogExit, "Exit status %d.", status)
		database.Close()
		os.Exit(status)
		} ()	
	
	flag.BoolVar(&flagReverseSort, "r", false, "Reverse output sort.")
	flag.BoolVar(&flagFollowOutput, "f", false, "Follow output.")
	flag.IntVar(&flagOutputLimit, "n", 0, "Number of lines to output. A zero value reports all lines.")
	flag.Parse()

	fmt.Println("Reverse:", flagReverseSort, "\t|")
	fmt.Println(" Follow:", flagFollowOutput, "\t|")
	fmt.Println("  Limit:", flagOutputLimit, "\t|")

	log(AWSLogStart, "Start %s.", strings.Trim(fmt.Sprint(os.Args), "[]"))

	switch rawCommand {
		case "help", "-h":
			status = printUsage()
						
		case "list":
			status = listBundles()
			
		case "log":
			status = printLog()
			
		case "history":
			status = printHistory()
			
		case "status":
			status = printStatus()
		
		case "report":
			status = printReport()
			
		default:
			fmt.Printf("Error: Unrecognized command '%v'.\n", rawCommand)
			status = 1
		}
	
	//	The end.
}

