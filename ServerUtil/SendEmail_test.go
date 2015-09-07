//  SendEmail  -  Send an Email using an existing account.
//
//  E.B.Smith  -  March, 2015


package ServerUtil


import (
    "testing"
)


//----------------------------------------------------------------------------------------
//                                                                           TestSendEmail
//----------------------------------------------------------------------------------------


func TestSendEmail(t *testing.T) {
    if true {
        error := SendEmail("smith.ed.b@gmail.com", "Hey Dude", "Message body.  Dude: This test worked.")
        if error != nil { t.Errorf("Got error sending email: %v.", error) }
    }
}

