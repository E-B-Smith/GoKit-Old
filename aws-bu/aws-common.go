package aws-bu

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

func log(logLevel AWSDebugLevel, processName string, message string) {
	_, error := db.Exec(
		"insert into AWSLogTable"
		"		(time, processname, level, message)"
		"values	(to_timestamp(?::float/1000000::float), '?', '?'::AWSLogLevel, '?');"
			time.Time, processName, logLevel, message)
}