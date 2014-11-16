//  invite-send-sms  -  Send an SMS to invite the  signed apps to the database.
//
//  E.B.Smith  -  November, 2014


package main


import (
	"os"
	"fmt"
	"strings"
	"io/ioutil"
	"net/url"
	"net/http"
	"encoding/base64"    
	)


var (
	twilioKey="AC6d456d2e0a49aad0d6ec1a92ddbd925f"
	twilioSecret="3191a48e7faa540b4c7b6e84c9a82402"
	twilioEncodedAuth=""
	twilioFromNumber="16072694306"
	twilioUrlString="https://api.twilio.com/2010-04-01/Accounts/"+twilioKey+"/Messages"
	)


func SendSMS(toNumber string, message string) error {
	//	Send a Twilio message to a phone number.

	if twilioEncodedAuth == "" {
		twilioEncodedAuth = "Basic "+base64.StdEncoding.EncodeToString([]byte(twilioKey+":"+twilioSecret))
		}

	formValues := url.Values{}
	formValues.Set("From", twilioFromNumber)
	formValues.Set("To", toNumber)
	formValues.Set("Body", message)

	request, _ := http.NewRequest("POST", twilioUrlString, strings.NewReader(formValues.Encode()))	
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
 	request.Header.Add("Authorization", twilioEncodedAuth)
	ZLog(ZLogDebug, "Request:\n%v", request)

	client := &http.Client{}
    response, error := client.Do(request)
    if error != nil {
        ZLog(ZLogError, "Can't contact Twilio: %v.", error)
    	}
    defer response.Body.Close()
    body, _ := ioutil.ReadAll(response.Body)

	ZLog(ZLogDebug, "Response:\n%v", response)
	ZLog(ZLogDebug, " Headers:\n%v", response.Header)
	ZLog(ZLogDebug, "    Body:\n%v", string(body))

	return error
	}


/*
create type UserStatus as enum
    (
     'UserStatusUnknown'             = 0
    ,'UserStatusPending'             = 1
    ,'UserStatusAccepted'            = 2
    ,'UserStatusSendInvite'          = 3
    ,'UserStatusInvited'             = 4
    ,'UserStatusDownloadFailed'      = 5
    ,'UserStatusDownloaded'          = 6
    ,'UserStatusFirstRun'            = 7
    ,'UserStatusActive'              = 8
    ,'UserStatusBlocked'             = 9
    ,'UserStatusSessionInvalid'      = 10
    );
*/


func main() {
	//	Send an email to each user with status 'UserStatusSendInvite' -- 

	error := connectDatabase()
	if error != nil { os.Exit(1) }
	defer disconnectDatabase()

	smsMessage := "Relcy beta has arrived.\nDownload the beta at\nhttps://relcy.com/invite?claim=%s"
		//	" ".   Send feedback at " "

	queryString := "select userID, phone, name, linkHashFromUserID(userID) from UserAnalyticsTable where userStatus = 3 and phone is not NULL;"
	rows, error := globalDatabase.Query(queryString)

	if error != nil {
		ZLog(ZLogDebug, "Can't read rows.\n%v", error)
		os.Exit(1)
		}
	defer rows.Close()

	var (count int; errorCount int)
	for rows.Next() {
		var (userID string; phone string; name string; hash string)
		error := rows.Scan(&userID, &phone, &name, &hash)
		if error != nil {
			errorCount++
			continue
		}
		ZLog(ZLogDebug, "Sending SMS to %s.", name)
		error = SendSMS(phone, fmt.Sprintf(smsMessage, hash))
		if error != nil {
			errorCount++
			ZLog(ZLogDebug, "Can't send SMS to %s.\n%v", error)
			continue
			}
		count++
		updateString := "update UserAnalyticsTable set userStatus=4 where userID=?;"
		_, error = globalDatabase.Exec(updateString, userID)
		if error != nil { ZLog(ZLogWarning, "User %s status not updated.\n%v", userID, error) }
		}
	ZLog(ZLogInfo, "Sent %d messages, %d errors.", count, errorCount)
	}


