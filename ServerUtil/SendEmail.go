//  SendEmail  -  Send an Email using an existing account.
//
//  E.B.Smith  -  March, 2015


package ServerUtil


import (
    "net/smtp"
    "../log"
)


func SendEmail(toAddress string, subject string, message string) error {
    //  Send an email --
    log.LogFunctionName()

    var (
        fromAddress     = "beinghappy@beinghappy.io"
        account         = "beinghappy@beinghappy.io"
        password        = "happy88!@"
        emailHost       = "smtp.gmail.com"
        smtpServer      = "smtp.gmail.com:587"
    )

    var error error
    toAddress, error = ValidatedEmailAddress(toAddress)
    if error != nil { return error }

    toArray := []string{toAddress}
    message = "Subject: "+subject+"\n\n"+message
    auth := smtp.PlainAuth("", account, password, emailHost)
    error = smtp.SendMail(smtpServer, auth, fromAddress, toArray, []byte(message))
    return error
}

