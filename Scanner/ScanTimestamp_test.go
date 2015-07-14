//  Strings.go  -  String functions.
//
//  E.B.Smith  -  November, 2014


package Scanner


import (
    "fmt"
    "time"
    "testing"
)


//----------------------------------------------------------------------------------------
//                                                                       TestScanTimestamp
//----------------------------------------------------------------------------------------


const kTestScanTimestampString =
`
0 2006-01-02T15:04:00-08:00
1 Jan 2 2006, 15:04
2 Jan 2 2006, 3:04PM
3 Monday, 02-Jan-06 16:04:00 EST
4 02 Jan 06 16:04 MST
5 02 Jan 06 15:04 -0800
6 Mon Jan 2 15:04:00 2006
7 Mon Jan 2 16:04:00 MST 2006
8 Mon Jan 02 15:04:00 -0800 2006
9 Mon, 02 Jan 2006 16:04:00 MST
10 Mon, 02 Jan 2006 15:04:00 -0800

11             this is an error string
`


func TestScanTimestamp(t *testing.T) {

    var error error
    var ts time.Time
    fmt.Printf("%s", kTestScanTimestampString)
//  utc, _ := time.LoadLocation("")
    testTime, error := time.Parse("2006-01-02T15:04:05Z07:00", "2006-01-02T15:04:00-08:00")
    if error != nil { panic(error); }
    scanner := NewStringScanner(kTestScanTimestampString)

    for i := 0; i < 11; i++ {
        scanner.ScanInt()
        ts, error = scanner.ScanTimestamp()
        if error != nil {
            //fmt.Printf("Test %d: Error %v.\n", i, error)
            t.Errorf("Test %d: Error %v.", i, error)
        } else if ! ts.Equal(testTime) {
            t.Errorf("Test %d: Scanned time %s but wanted %s. Diff: %v\nInput: %s",
                i, ts.Format(time.RFC3339), testTime.Format(time.RFC3339),
                ts.Sub(testTime), scanner.Token())
        }
        //fmt.Printf("\n\n\n\n")
    }

    scanner.ScanInt()
    ts, error = scanner.ScanTimestamp()
    if error == nil {
        t.Errorf("Expected an error but got %v.", ts)
    }
}

