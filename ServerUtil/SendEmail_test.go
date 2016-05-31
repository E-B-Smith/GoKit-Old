

//----------------------------------------------------------------------------------------
//
//                                                                       SendEmail_test.go
//                                                  ServerUtil: Basic API server utilities
//
//                                                                   E.B.Smith, March 2015
//                        -©- Copyright © 2015-2016 Edward Smith, all rights reserved. -©-
//
//----------------------------------------------------------------------------------------


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

