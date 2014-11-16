//  invite-add-apps  -  Adds signed apps to the database.
//
//  E.B.Smith  -  November, 2014


package main


import (
	"os"
	"bufio"
	"strings"
    "strconv"
	)


func insertItemsInFile(file *os.File) {
	var linecount = 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		linecount++
		var platform = 0
		fields := strings.Fields(scanner.Text())
		if len(fields) == 0 {
			continue
			}
		if len(fields) >= 2 {
			platform, _ = strconv.Atoi(fields[1])
			}
		if platform < 1 || platform > 2 {
			ZLog(ZLogError, "Line %d. Expected 'InviteID'<white-space>'platformID' on the same line.", linecount)
			os.Exit(1)
			}
		
		_, error := globalDatabase.Exec("insert into AvailableInviteTable (inviteID, platformType) values (?, ?);", fields[0], platform)
		if error != nil {
			ZLog(ZLogError, "Line %d. Can't inserting record.\n%v", error)
			os.Exit(1)
			}
		}
	}


func main() {
	//	Process files on the command line or stdin -- 

	error := connectDatabase()
	if error != nil {
		os.Exit(1)
	}
	defer disconnectDatabase()

	if len(os.Args) <= 1 {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			ZLog(ZLogDebug, "Reading from stdin.")
			insertItemsInFile(os.Stdin)
		} else {
			ZLog(ZLogError, "Can't read stdin.")
		}
	} else  {
		for i:= 0; i < len(os.Args); i++ {
			ZLog(ZLogDebug, "Reading '%s'...", os.Args[i])
			var file *os.File = nil
			file, error = os.Open(os.Args[i])
			if error != nil {
	            ZLog(ZLogError, "Error: Can't open file '%s' for reading: %v.", os.Args[i], error)
	            os.Exit(1)
	        	}
			defer file.Close()
			insertItemsInFile(file)
			}
		}
	}

