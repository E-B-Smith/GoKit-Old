//  Scanner  -  A parser to parse a configuration file.
//
//  E.B.Smith  -  November, 2014


package Scanner


import (
    "time"
    "strings"
)


func (scanner *Scanner) ScanTimestamp() (time.Time, error) {

    local, error := time.LoadLocation("Local")
    if error != nil { panic(error); }

    utc, error := time.LoadLocation("")
    if error != nil { panic(error); }

    type TimestampFormat struct {
        Location    *time.Location
        Format      string
    }
    kTimestampFormats := []TimestampFormat{
        {utc,    "2006-01-02T15:04:05Z07:00"},        //  1
        {local,  "Jan 2 2006, 15:04"},                //  4
        {local,  "Jan 2 2006, 3:04PM"},               //  4
        {utc,    "Monday, 02-Jan-06 15:04:05 MST"},   //  4
        {utc,    "02 Jan 06 15:04 MST"},          //  5
        {utc,    "02 Jan 06 15:04 -0700"},        //  5
        {local,  "Mon Jan 2 15:04:05 2006"},      //  5
        {utc,    "Mon Jan 2 15:04:05 MST 2006"},      //  6
        {utc,    "Mon Jan 02 15:04:05 -0700 2006"},   //  6
        {utc,    "Mon, 02 Jan 2006 15:04:05 MST"},    //  6
        {utc,    "Mon, 02 Jan 2006 15:04:05 -0700"},  //  6
    }

    stringInput := ""
    stringParts := 0
    formatParts := 0
    var resultTime time.Time
    for _, kFormat := range(kTimestampFormats) {
        scanner.error = nil

        //  Count the parts --

        formatParts = 0
        parts := strings.Split(kFormat.Format, " ")
        for _, part := range(parts) {
            if part != "" { formatParts++; }
        }

        //  Get the parts --

        for stringParts < formatParts {
            var s string
            if stringInput != "" { stringInput += " " }
            scanner.error = scanner.ScanSpaces()
            if scanner.error != nil { return time.Time{}, scanner.error }
            s, scanner.error = scanner.ScanString()
            if scanner.error != nil { return time.Time{}, scanner.error }
            stringInput += s
            stringParts++

            var r rune
            r, _, _ = scanner.reader.ReadRune()
            if r == ',' || r == ';' {
                stringInput += string(r)
            } else {
                scanner.reader.UnreadRune()
            }
        }

        //  Try to parse --

        resultTime, scanner.error = time.ParseInLocation(kFormat.Format, stringInput, kFormat.Location)
            //fmt.Printf("Format parts: %d String parts: %d Input: '%s' Format: '%s' Error: %v.\n",
            //  formatParts, stringParts, stringInput, kFormat.Format, scanner.error)
        if scanner.error == nil {
            scanner.token = stringInput
            return resultTime, scanner.error
        }
    }

    return resultTime, scanner.error
}


