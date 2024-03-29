

//----------------------------------------------------------------------------------------
//
//                                                                            SendEmail.go
//                                                  ServerUtil: Basic API server utilities
//
//                                                                   E.B.Smith, March 2015
//                        -©- Copyright © 2015-2016 Edward Smith, all rights reserved. -©-
//
//----------------------------------------------------------------------------------------


package ServerUtil


import (
    "io"
    "fmt"
    "net"
    "errors"
    "regexp"
    "strings"
    "net/smtp"
    "crypto/tls"
    "violent.blue/GoKit/Log"
    "violent.blue/GoKit/Util"
)


//  Send an email --
func (config Configuration) SendEmail(toAddress string, subject string, message string) (err_r error) {
    Log.LogFunctionName()

    /*
    Uses:
        config.EmailFromAddress     "Blitz <blitz@blitzhere.com>"
        config.EmailAccount         blitz@blitzhere.com
        config.EmailPassword        *****
        config.EmailSMPTServer      smtp.gmail.com:587
    */

    defer func() {
        if r := recover(); r != nil {
            switch x := r.(type) {
            case string:
                err_r = errors.New(x)
            case error:
                err_r = x
            default:
                err_r = errors.New("Unknown panic")
            }
        }
    } ()

    checkError := func(error error) {
        if error != nil {
            Log.Errorf("Panic with error: %+v.", error)
            Log.LogStackWithError(error)
            panic(error)
        }
    }

    //  Get the host name:

    emailHostArray := strings.Split(config.EmailSMTPServer, ":")
    if len(emailHostArray) == 0 {
        return fmt.Errorf("No email host configured.")
    }
    SMTPServer := emailHostArray[0]

    //  Get the raw 'from' address --

    var error error
    var rexp *regexp.Regexp
    rexp, error = regexp.Compile("<(.*?)>")
    checkError(error)

    rawFromAddress := config.EmailFromAddress
    fromArray := rexp.FindAllString(rawFromAddress, -1)
    if len(fromArray) > 0 {
        rawFromAddress = fromArray[0]
        rawFromAddress = strings.Trim(rawFromAddress, " <>")
    }
    Log.Debugf("Full from address: '%s' Raw address: '%s'.", config.EmailFromAddress, rawFromAddress)

    toAddress, error = Util.ValidatedEmailAddress(toAddress)
    checkError(error)

    //  Now connect with the server using TLS --

    tlsconfig := &tls.Config {
        // InsecureSkipVerify: true,   //  eDebug -- Remove for production.
        ServerName:         SMTPServer,
    }

    var client *smtp.Client
    var tlsDial bool = false

    //  Connect --

    if tlsDial {

        var connection *tls.Conn
        Log.Debugf("TLS Connection to '%s'.", config.EmailSMTPServer)
        connection, error = tls.Dial("tcp", config.EmailSMTPServer, tlsconfig)
        checkError(error)

        client, error = smtp.NewClient(connection, SMTPServer)
        checkError(error)

    } else {

        var connection net.Conn
        Log.Debugf("Connecting to '%s'.", config.EmailSMTPServer)
        connection, error = net.Dial("tcp", config.EmailSMTPServer)
        checkError(error)

        client, error = smtp.NewClient(connection, SMTPServer)
        checkError(error)

        Log.Debugf("Starting TLS.")
        error = client.StartTLS(tlsconfig)
        checkError(error)

    }

    defer func(client *smtp.Client) {
        error := client.Quit()
        if error != nil {
            if strings.HasPrefix(error.Error(), "250 2.0.0 OK") {
                error = nil
                Log.Debugf("Email sent.")
            } else {
                Log.LogError(error)
            }
        } else {
            Log.Debugf("Email sent.")
        }
    } (client)

    //  Auth & send email --

    Log.Debugf("Auth '%s' / '%s'", config.EmailAccount, config.EmailPassword)
    auth := smtp.PlainAuth("", config.EmailAccount, config.EmailPassword, SMTPServer)
    Log.Debugf("Send account: '%s' send server: '%s' To: '%s'.",
        config.EmailAccount,
        SMTPServer,
        toAddress,
    )

    error = client.Auth(auth)
    checkError(error)

    Log.Debugf("From: %s", rawFromAddress)
    error = client.Mail(rawFromAddress)
    checkError(error)

    error = client.Rcpt(toAddress)
    checkError(error)

    var writer io.Writer
    writer, error = client.Data()
    checkError(error)

    message = fmt.Sprintf("Subject: %s\nFrom: %s\n\n%s\n",
        subject,
        config.EmailFromAddress,
        message,
    )
    //Log.Debugf("Message: '%s'.", message)
    _, error = writer.Write([]byte(message))
    checkError(error)

    Log.Debugf("Email sent.  Confirming...")
    return error
}

