//  pgsql.array  -  A go Postgres interface for handling postgres arrays.
//
//  E.B.Smith  -  November, 2014


package pgsql


import (
    "fmt"
    "time"
    "strings"
    "strconv"
    "database/sql"
    "violent.blue/GoKit/Log"
    "violent.blue/GoKit/Scanner"
)


//----------------------------------------------------------------------------------------
//
//                                                                                  Arrays
//
//----------------------------------------------------------------------------------------


//----------------------------------------------------------------------------------------
//                                                                                 Strings
//----------------------------------------------------------------------------------------


func NullStringFromStringArray(ary []string) sql.NullString {
    if len(ary) == 0 {
        Log.Debugf("NULL returned for empty string.")
        return sql.NullString {}
    }

    var result string = "{\""+ary[0];
    for i:=1; i < len(ary); i++ {
        result += "\",\""+ary[i]
    }
    result += "\"}"

    Log.Debugf("Output: |%s|.", result)
    return sql.NullString { Valid: true, String: result }
}


/*
func StringArrayFromString(s *string) []string {
    if s == nil { return *new([]string) }

    str := strings.Trim(*s, "{}")
    a := make([]string, 0, 10)
    for _, ss := range strings.Split(str, ",") {
        a = append(a, ss)
    }
    return a
}
*/


func StringArrayFromNullString(nullstring sql.NullString) []string {

    array := make([]string, 0, 10)

    if ! nullstring.Valid {
        Log.Debugf("NULL returned for empty string.")
        return array
    }

    Log.Debugf("Input: |%s|.", nullstring.String)

    scanner := Scanner.NewStringScanner(strings.Trim(nullstring.String, "{} "))
    for ! scanner.IsAtEnd() {
        s, error := scanner.ScanSQLString();
        if error != nil { Log.LogError(error) }

        if len(s) > 0 {
            if s == "NULL" { s = "" }
            Log.Debugf("Adding %s", s)
            array = append(array, s)
        }

        scanner.ScanSpaces()
        c := scanner.NextRune()   //  Should be a comma or nothing.
        Log.Debugf("Sep is '%c' At end: %t Token: '%s'.", c, scanner.IsAtEnd(), scanner.Token())
        if ! (c == ',' || c == 0) {
            panic(fmt.Errorf("Mal-formed postgres string array. Input was: '%s'. Error at '%c'.",
                nullstring.String, c))
        }
    }

    Log.Debugf("Output: %v.", array)

    return array
}


//----------------------------------------------------------------------------------------
//                                                                                   Int32
//----------------------------------------------------------------------------------------


func StringFromInt32Array(ary []int32) string {
    if len(ary) == 0 {
        return "{}"
    }

    var result string = "{"+strconv.Itoa(int(ary[0]));
    for i:=1; i < len(ary); i++ {
        result += ","+strconv.Itoa(int(ary[i]))
    }
    result += "}"
    return result
}


func Int32ArrayFromString(s *string) []int32 {
    if s == nil { return *new([]int32) }

    str := strings.Trim(*s, "{}")
    a := make([]int32, 0, 10)
    for _, ss := range strings.Split(str, ",") {
        i, error := strconv.Atoi(ss)
        if error == nil { a = append(a, int32(i)) }
    }
    return a
}


//----------------------------------------------------------------------------------------
//                                                                                 Float64
//----------------------------------------------------------------------------------------


func Float64ArrayFromNullString(s *sql.NullString) []float64 {
    if s == nil || !s.Valid {
        return *new([]float64)
    }

    a := make([]float64, 0, 10)
    str := strings.Trim(s.String, "{}")
    for _, ss := range strings.Split(str, ",") {
        f, error := strconv.ParseFloat(ss, 64)
        if error == nil { a = append(a, f) }
    }
    return a
}


func StringFromFloat64Array(ary []float64) string {
    if len(ary) == 0 {
        return "{}"
    }

    var result string = "{"+strconv.FormatFloat(ary[0], 'g', -1, 64);
    for i:=1; i < len(ary); i++ {
        result += ","+strconv.FormatFloat(ary[i], 'g', -1, 64)
    }
    result += "}"
    return result
}


//----------------------------------------------------------------------------------------
//                                                                                    time
//----------------------------------------------------------------------------------------


func TimeArrayFromNullString(s *sql.NullString) []time.Time {
    if s == nil || !s.Valid {
        return *new([]time.Time)
    }

    const kFormat = "2006-01-02 15:04:05-07" // "2016-08-08 22:30:00+00" <>  "2006-01-02T15:04:05Z07:00"
    Log.Debugf("Parsing '%s'.", s.String)

    a := make([]time.Time, 0, 10)
    str := strings.Trim(s.String, "{}")
    for _, ss := range strings.Split(str, ",") {
        ss = strings.Trim(ss, " \"")
        t, error := time.Parse(kFormat, ss)
        if error == nil { a = append(a, t) }
    }

    Log.Debugf("Returned %d items.", len(a))
    return a
}


func NullStringFromTimeArray(ary []time.Time) sql.NullString {
    if len(ary) == 0 {
        return sql.NullString {}
    }

    var result string = "{"+ary[0].Format(time.RFC3339Nano)
    for i:=1; i < len(ary); i++ {
        result += ","+ary[i].Format(time.RFC3339Nano)
    }
    result += "}"
    return sql.NullString { Valid: true, String: result }
}

