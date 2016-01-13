//  Strings.go  -  String functions.
//
//  E.B.Smith  -  November, 2014


package Util


import (
    "fmt"
    "time"
    "testing"
    )


//----------------------------------------------------------------------------------------
//                                                      TestStringIncludingCharactersInSet
//----------------------------------------------------------------------------------------


func TestStringIncludingCharactersInSet(t *testing.T) {
    tests := []struct {
        testin, testout string
    } {
        { "123-#-456",          "123456" },
        { "123456",             "123456" },
        { "",                   "" },
        { "-123456-",           "123456" },
        { "aslkdjlaskj",        "" },
    }

    for _, test := range tests {
        result := StringIncludingCharactersInSet(test.testin, "1234567890")
        if false { fmt.Printf("%s\t\t%s\n", result, test.testout) }
        if result != test.testout {
            t.Errorf("Got %s but want %s.", result, test.testout)
        }
    }
}



//----------------------------------------------------------------------------------------
//                                                                          TestHumanBytes
//----------------------------------------------------------------------------------------


func TestHumanBytes(t *testing.T) {
    tests := []struct {
        testin int64; testout string
    } {
        { 123,          "123 B" },
        { 1234,         "1.21 KB" },
        { 123456,       "120.56 KB" },
        { 1234567,      "1.18 MB" },
        { 123456789,    "117.74 MB" },
        { 1234567890,   "1.15 GB" },
      }

    for _, test := range tests {
        result := HumanBytes(test.testin)
        if false { fmt.Printf("%d\t%s\t%s\n", test.testin, result, test.testout) }
        if result != test.testout {
            t.Errorf("Got %s but want %s.", result, test.testout)
        }
    }
}



//----------------------------------------------------------------------------------------
//                                                                            TestHumanInt
//----------------------------------------------------------------------------------------


func TestHumanInt(t *testing.T) {
    tests := []struct {
        testin int64; testout string
    }{
        { 0,            "0" },
        { -0,           "0" },
        { 123,          "123" },
        { 1234,         "1,234" },
        { -1234,        "-1,234" },
        { 123456,       "123,456" },
        { 1234567,      "1,234,567" },
        { 123456789,    "123,456,789" },
        { -123456789,   "-123,456,789" },
    }

    for _, test := range tests {
        result := HumanInt(test.testin)
        if false { fmt.Printf("%d\t%s\t%s\n", test.testin, result, test.testout) }
        if result != test.testout {
            t.Errorf("Got %s but want %s.", result, test.testout)
        }
    }
}

//----------------------------------------------------------------------------------------
//                                                               TestCompareVersionStrings
//----------------------------------------------------------------------------------------


func TestCompareVersionStrings(t *testing.T) {
    const NSOrderedSame = 0
    const NSOrderedAscending = -1
    const NSOrderedDescending = 1

    tests := []struct {
        Result int; Version1 string; Version2 string
    }{
        { NSOrderedSame,        "1.2.1",        "1.2.1" },
        { NSOrderedAscending,   "1.2.1",        "1.2.2" },
        { NSOrderedDescending,  "1.2.2",        "1.2.1" },
        { NSOrderedSame,        "1.2.1.3",      "1.2.1.3" },
        { NSOrderedAscending,   "1.2.1",        "1.2.1.3" },

        { NSOrderedSame,        "1.020.1",      "1.20.1" },
        { NSOrderedAscending,   "1.02.1",       "1.020.2" },
        { NSOrderedDescending,  "3.2.2",        "2.3.2" },

        { NSOrderedSame,        "1.2.00.00",    "1.2" },
        { NSOrderedAscending,   "1.2",          "1.2.0.1" },
        { NSOrderedDescending,  "1.2.0.1",      "1.2" },
    }

    for _, test := range tests {
        result := CompareVersionStrings(test.Version1, test.Version2)
        if false { fmt.Printf("%s\t%s\t%d: %d\n", test.Version1, test.Version2, test.Result, result) }
        if result != test.Result {
            t.Errorf("%s\t%s\t should be %d but got %d\n", test.Version1, test.Version2, test.Result, result)
        }
    }
}

//----------------------------------------------------------------------------------------
//                                                                       TestHumanDuration
//----------------------------------------------------------------------------------------


func TestHumanDuration(t *testing.T) {
    var d time.Duration =
        time.Hour * 24 +
        time.Hour * 2 +
        time.Minute * 3 +
        time.Second * 4 +
        time.Second / 2
    s := HumanDuration(d)
    var e string = "1 day 2:03:04.5 hours"
    if s != e {
         t.Errorf("Expected '%s' but got '%s'.\n", e, s)
    }
}



