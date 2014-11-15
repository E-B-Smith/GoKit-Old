//  invite-send-sms  -  Send an SMS to invite the  signed apps to the database.
//
//  E.B.Smith  -  November, 2014


package main


import (
	"os"
	"bufio"
	"strings"
    "strconv"
	)

// function percentEscapeString()
// 	{
// 	perl -lpe 's/([^A-Za-z0-9])/sprintf("%%%02X", ord($1))/seg' <<<"$1"
// 	}

// twiliokey=AC6d456d2e0a49aad0d6ec1a92ddbd925f
// twiliosecret=3191a48e7faa540b4c7b6e84c9a82402
// twiliofrom=16072694306
// twilioto=6503916685   # Rohit
// #twilioto=2404856540   # Sap
// #twilioto=4086379251  # Hemant
// #twilioto=4156152570  # Edward
// encodedauth=$(echo "$twiliokey:$twiliosecret" | base64)
// body="McD is still better."
// body=$(percentEscapeString "$body")

// curl -X POST \
//     --insecure --retry 3 --silent --show-error \
//     -H "Authorization: Basic $encodedauth" \
//     --data "From=$twiliofrom&To=$twilioto&Body=$body" \
//         "https://api.twilio.com/2010-04-01/Accounts/$twiliokey/Messages"

// echo ""

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
			ZLog(ZLogError, "Can't open stdin.")
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

