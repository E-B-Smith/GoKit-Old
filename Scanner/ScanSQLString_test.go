//  Strings.go  -  String functions.
//
//  E.B.Smith  -  November, 2014


package Scanner


import (
    "testing"
)


//----------------------------------------------------------------------------------------
//                                                                       TestScanTimestamp
//----------------------------------------------------------------------------------------


type TestCaseType struct {
    Test    string
    Result  string
    Success bool
}


func TestScanSQLString(t *testing.T) {

    TestCases := []TestCaseType {
        { "Unquoted  ",       "Unquoted",     true  },
        { "'Quoted'  ",       "Quoted",       true  },
        { "'Don''t look!'  ", "Don't look!",  true  },
        { "",                 "",             false },
        { "''",               "",             true  },
        { "'fail",            "",             false },
    }

    for _, tc := range TestCases {

        scanner := NewStringScanner(tc.Test)
        result, error := scanner.ScanSQLString()
        t.Logf("Tested '%s' expected '%s' got '%s' error: %v.", tc.Test, tc.Result, result, error)
        if error != nil {
            if tc.Success {
                t.Errorf("Failed!")
            }
        }
        if result != tc.Result {
            t.Errorf("Failed!")
        }
    }
}

