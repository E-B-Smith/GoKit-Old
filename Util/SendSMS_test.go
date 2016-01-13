//  SendSMS  -  Send an SMS to a phone using Twilio.
//
//  E.B.Smith  -  March, 2015


package Util


import (
    "testing"
)


//----------------------------------------------------------------------------------------
//                                                                             TestSendSMS
//----------------------------------------------------------------------------------------


func TestSendSMS(t *testing.T) {
    if false {
        error := SendSMS("4156152570", "Dude: Test worked.")
        if error != nil { t.Errorf("Got error sending email: %v.", error) }
    }
}

