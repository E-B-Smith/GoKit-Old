//  claim-invite  -  Claim an invite link.
//
//  E.B.Smith  -  November, 2014


package main


import (
    "fmt"
    "net/http"
    "strconv"
    "strings"
	)


var globalPort = 8080;


func sendErrorPage(writer http.ResponseWriter, message string) {
	fmt.Fprintf(writer, message)
	}

func sendGenericError(writer http.ResponseWriter) {
	sendErrorPage(writer, "<p>Sorry, your invite can not be claimed at this time.<br><br>Try again later.</p>")
	}


func claimInvite(writer http.ResponseWriter, request *http.Request) {
	//	Redirect to the app download page -- 

	var platform int
	var redirectString string 

	browser := request.Header.Get("User-Agent")
	browser = strings.ToLower(browser)
	if strings.Contains(browser, "iphone") {
		platform = 1
		redirectString = "itms-services://?action=download-manifest&url=https://relcy.com/alpha/%s/SearchStaging.plist"
	} else {
	if strings.Contains(browser, "android") || strings.Contains(browser, "linux") {
		platform = 2
		redirectString = "https://www.relcy.com/alpha/%s/app_beta.apk";
	} else {
		ZLog(ZLogWarning, "Page opened in non-device browser.  Header: %v.", request.Header)
		sendErrorPage(writer, "To download Relcy Look, open this link on your iPhone or Android device.")
		return
		}}

	inviteCode := request.URL.Query().Get("claim")
	ZLog(ZLogDebug, "Handle path %s claim %s platform %d.", request.URL.Path, inviteCode, platform)
	row := globalDatabase.QueryRow("select inviteIDFromLinkHashAndPlatform(?, ?);", inviteCode, 2)

	var inviteID string
	error := row.Scan(&inviteID)
	if error != nil || inviteID == "" {
		ZLog(ZLogError, "Can't get inviteID for invite code %s:\n%v\n.", inviteCode, error, row);
		sendGenericError(writer)
		return
		}

	urlString := fmt.Sprintf(redirectString, inviteID)

	ZLog(ZLogDebug, "Got invite ID '%s'.", inviteID)
	ZLog(ZLogDebug, "Redirecting to '%s'.", urlString)
	http.Redirect(writer, request, urlString, 307)
	sendErrorPage(writer, "Thank you for trying Relcy Look.")
	}


func main() {
	connectDatabase();
	defer disconnectDatabase();
    http.HandleFunc("/invite", claimInvite)
    http.HandleFunc("/invite/", claimInvite)
    http.ListenAndServe(":"+strconv.Itoa(globalPort), nil)
	}

