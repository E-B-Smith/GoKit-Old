

//----------------------------------------------------------------------------------------
//
//                                                                    TemplateFunctions.go
//                                                     ServerUtil: HTTP Template Functions
//
//                                                                   E.B.Smith, March 2016
//                        -©- Copyright © 2015-2016 Edward Smith, all rights reserved. -©-
//
//----------------------------------------------------------------------------------------


package ServerUtil


import (
    "fmt"
    "math"
    "time"
    "html"
    "violent.blue/GoKit/Log"
    "violent.blue/GoKit/Util"
)


//----------------------------------------------------------------------------------------
//                                                                 HTTP Template Functions
//----------------------------------------------------------------------------------------


//  Un-escape an HTML string in a template
func unescapeHTMLString(args ...interface{}) string {
    Log.LogFunctionName()
    Log.Debugf("%+v", args...)
    ok := false
    var s string
    if len(args) == 1 {
        s, ok = args[0].(string)
        s = html.UnescapeString(s)
    }
    if !ok {
        s = fmt.Sprint(args...)
    }
    return s
}


//  Escape an HTML string in a template
func escapeHTMLString(args ...interface{}) string {
    Log.LogFunctionName()
    Log.Debugf("%+v", args...)
    ok := false
    var s string
    if len(args) == 1 {
        s, ok = args[0].(string)
        s = html.EscapeString(s)
    }
    if !ok {
        s = fmt.Sprint(args...)
    }
    return s
}


//  Evaluate a boolean pointer in a template
func boolPtr(b *bool) bool {
    if b != nil && *b { return true }
    return false
}


func stringPtr(s *string) string {
    if s == nil { return "" }
    return *s
}


const kMonthYearFormat string = "1/2006"


func timeFromDouble(timestamp float64) time.Time {
    i, f := math.Modf(timestamp)
    var  sec int64 = int64(math.Floor(i))
    var nsec int64 = int64(f * 1000000.0)
    return time.Unix(sec, nsec)
}

//  Format a date ptr.
func MonthYearStringFromEpochPtr(epoch *float64) string {
    if epoch == nil || *epoch <= 0.0 {
        return ""
    }
    t := timeFromDouble(*epoch)
    return t.Format(kMonthYearFormat)
}


//  Parse a date format:
func ParseMonthYearString(s string) time.Time {
    s = Util.StringIncludingCharactersInSet(s, "0123456789/")
    t, _ := time.Parse(kMonthYearFormat, s)
    return t;
}

