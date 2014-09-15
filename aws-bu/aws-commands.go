package main

import (
//	"os"
	"fmt"
	"strconv"
//	"flag"
	"time"
	"bytes"
)


func listBundles() int {
	rows, err := database.Query("select bundle from awsobjecttabletotals;")
	if err != nil {
		log(AWSLogError, "Database error: %v.", err)
		return 1
	}
	var bundleName string
	for rows.Next() {
		rows.Scan(&bundleName)
		fmt.Println(bundleName)
	}
	return 0
}


func printLog() int {
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
	
	rows, err := database.Query(queryBuffer.String())
	if err != nil {
		log(AWSLogError, "Database error: %v.", err)
		return 1;
	}
	
	var timestamp time.Time
	var processname string
	var level string
	var message string
	for rows.Next() {
		rows.Scan(&timestamp, &processname, &level, &message)
		fmt.Printf("%s  %-12s %-12s: %s\n", timestamp, processname, level, message)
	}
	
	return 0
}


func printHistory() int {
	var queryBuffer bytes.Buffer
	queryBuffer.WriteString("select time, processname, level, message from awslogtable ")
	queryBuffer.WriteString("where level='AWSLogStart'::AWSLogLevel or level='AWSLogExit'::AWSLogLevel ")

	if flagReverseSort {
		queryBuffer.WriteString("order by entry desc ")
	}
	
	if flagOutputLimit != 0 {
		queryBuffer.WriteString("limit ")
		queryBuffer.WriteString(strconv.Itoa(flagOutputLimit))
	}
	
	queryBuffer.WriteString(";")
	
	log(AWSLogDebug, "History query string: %s.", queryBuffer.String())

	rows, err := database.Query(queryBuffer.String())
	if err != nil {
		log(AWSLogError, "Database error: %v.", err)
		return 1;
	}
	
	var timestamp time.Time
	var processname string
	var level string
	var message string
	for rows.Next() {
		rows.Scan(&timestamp, &processname, &level, &message)
		fmt.Printf("%s  %-12s %-12s: %s\n", timestamp, processname, level, message)
	}
	
	return 0
}


func printStatus() int {
	fmt.Printf("printStatus");
	return 0
}

