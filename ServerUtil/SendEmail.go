//  SendEmail  -  Send an Email using an existing account.
//
//  E.B.Smith  -  March, 2015


package ServerUtil


import (
    "fmt"
    "regexp"
    "strings"
    "net/smtp"
    "violent.blue/GoKit/Log"
    "violent.blue/GoKit/Util"
)


func (config Configuration) SendEmail(toAddress string, subject string, message string) error {
    //  Send an email --
    Log.LogFunctionName()

/*
    Uses:
        config.EmailAccount
        config.Password
        config.EmailSMTPServer
*/

    //  Get the host name:

    emailHostArray := strings.Split(config.EmailSMTPServer, ":")
    if len(emailHostArray) == 0 {
        return fmt.Errorf("No email host configured.")
    }
    emailHost := emailHostArray[0]

    //  Get the raw 'from' address --

    var error error
    var rexp *regexp.Regexp
    rexp, error = regexp.Compile("<(.*?)>")
    if error != nil { Log.LogError(error); return error; }

    rawFromAddress := config.EmailAddress
    fromArray := rexp.FindAllString(rawFromAddress, -1)
    if len(fromArray) > 0 {
        rawFromAddress = fromArray[0]
        rawFromAddress = strings.Trim(rawFromAddress, " <>")
    }
    Log.Debugf("From: '%s' '%s'.", config.EmailAddress, rawFromAddress)

    toAddress, error = Util.ValidatedEmailAddress(toAddress)
    if error != nil { return error }

    toArray := []string{toAddress}
    message = fmt.Sprintf("Subject: %s\nFrom: %s\n\n%s\n", subject, config.EmailAddress, message)
    auth := smtp.PlainAuth("", config.EmailAccount, config.EmailPassword, emailHost)
    error = smtp.SendMail(config.EmailSMTPServer, auth, rawFromAddress, toArray, []byte(message))
    return error
}

