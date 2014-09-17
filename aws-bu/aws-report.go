package main

import (
	"fmt"
	"bytes"
)


func printReport() int {
	var queryBuffer bytes.Buffer
	queryBuffer.WriteString("select * from awsobjecttabletotals;")
	
	rows, error := database.Query(queryBuffer.String())
	if error != nil {
		log(AWSLogError, "Database error: %v.", error)
		return 1;
	}

	fmt.Println(rows);

	// var timestamp time.Time
	// var processname string
	// var level string
	// var message string
	// for rows.Next() {
	// 	rows.Scan(&timestamp, &processname, &level, &message)
	// 	fmt.Printf("%s  %-12s %-12s: %s\n", timestamp, processname, level, message)
	// }
	
	return 0
}
