package main

import (
	"os"
	"fmt"
	"strconv"	
	"flag"
	"time"
	"bytes"
	_ "github.com/lib/pq"
	"database/sql"
)


var (
	flagReverseSort bool = false
	flagFollowOutput bool = false
	flagOutputLimit int = 0
)


func main() {

	fmt.Println("Arguments:", os.Args)

	if len(os.Args) <= 1 {
		fmt.Println("Error: Run command expected.")
		printUsage()
		os.Exit(1)
	}

	var command string = os.Args[1]
	fmt.Println("Command:", command)
	os.Args = os.Args[1:]

	flag.BoolVar(&flagReverseSort, "r", false, "Reverse output sort.")
	flag.BoolVar(&flagFollowOutput, "f", false, "Follow output.")
	flag.IntVar(&flagOutputLimit, "n", 0, "Number of lines to output. Set to zero to output all lines.")
	flag.Parse()

	fmt.Println("Reverse:", flagReverseSort, "\t|")
	fmt.Println(" Follow:", flagFollowOutput, "\t|")
	fmt.Println("  Limit:", flagOutputLimit, "\t|")

	switch command {
		case "help", "-h":
			printUsage()
						
		case "bundles":
			listBundles()
			
		case "log":
			printLog()
			
		case "history":
			printHistory()
			
		case "status":
			printStatus()
		}
			
/*	if len(flag.Args()) != 0 {
		fmt.Println("Error: Invalid command line argument '"+flag.Arg(0)+"'.")
		os.Exit(1)
	}


	printUsage()
	listBundles()
	printLog()
*/
}


func printUsage() {
    fmt.Fprintf(os.Stderr, "usage: aws-bu [ help | bundles | log | history | status | report ] [ <options> ]\n\noptions:\n")
    flag.PrintDefaults()
}


func listBundles() {
	db := connectDatabase()
	rows, err := db.Query("select bundle from awsobjecttabletotals;")
	if err != nil {
		panic(err)
	}
	var bundleName string
	for rows.Next() {
		rows.Scan(&bundleName)
		fmt.Println(bundleName)
	}
}


func printLog() {
	db := connectDatabase()

	var queryBuffer bytes.Buffer
	queryBuffer.WriteString("select time, processname, level, message from awslogtable")
		
	if flagReverseSort {
		queryBuffer.WriteString(" order by entry desc")
	}
	
	if flagOutputLimit != 0 {
		queryBuffer.WriteString(" limit ")
		queryBuffer.WriteString(strconv.Itoa(flagOutputLimit))
	}
	
	queryBuffer.WriteString(";")
	
	rows, err := db.Query(queryBuffer.String())
	if err != nil {
		panic(err)
	}
	
	var timestamp time.Time
	var processname string
	var level string
	var message string
	for rows.Next() {
		rows.Scan(&timestamp, &processname, &level, &message)
		fmt.Printf("%s  %-12s %-12s: %s\n", timestamp, processname, level, message)
	}
}
	

func printHistory() {
	fmt.Printf("printHistory");
/*	db := connectDatabase()

	var queryBuffer bytes.Buffer
	queryBuffer.WriteString("select time, processname, level, message from awslogtable")
		
	if flagReverseSort {
		queryBuffer.WriteString(" order by entry desc")
	}
	
	if flagOutputLimit != 0 {
		queryBuffer.WriteString(" limit ")
		queryBuffer.WriteString(strconv.Itoa(flagOutputLimit))
	}
	
	queryBuffer.WriteString(";")
	
/*	rows, err := db.Query(queryBuffer.String())
	if err != nil {
		panic(err)
	}
*/}


func printStatus() {
	fmt.Printf("printStatus");
}


func connectDatabase() *sql.DB {
	conn, err := sql.Open("postgres", "user=Edward dbname=Edward sslmode=disable")
	if err != nil {
		panic(err)
	}
	return conn
}


