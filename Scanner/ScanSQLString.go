//  ScanSQLString  -  A parser to parse a configuration file.
//
//  E.B.Smith  -  November, 2014


package Scanner


import (
    "bytes"
    "errors"
    //"violent.blue/GoKit/Log"
)


func (scanner *Scanner) ScanSQLString() (string, error) {
    scanner.error = scanner.ScanSpaces()

    var r rune
    r, _, scanner.error = scanner.reader.ReadRune()
    if scanner.error != nil {
        return "", scanner.error
    }

    switch r {

    case '"':
        scanner.reader.UnreadRune()
        return scanner.ScanQuotedString()

    case '\'':
        break

    default:
        scanner.reader.UnreadRune()
        return scanner.ScanString()
    }

    var buffer bytes.Buffer
    for true {
        r, _, scanner.error = scanner.reader.ReadRune()
        if  scanner.error != nil {
            scanner.token = ""
            scanner.error = errors.New("quote error")
            return scanner.token, scanner.error
        }

        switch r {
        case '\'':
            r, _, scanner.error = scanner.reader.ReadRune()
            if scanner.error != nil {
                scanner.token = buffer.String()
                scanner.error = nil
                return scanner.token, scanner.error
            }
            if r == '\'' {
                buffer.WriteRune(r)
            } else {
                scanner.token = buffer.String()
                return scanner.token, nil
            }
        default:
            buffer.WriteRune(r)
        }
    }
    scanner.token = buffer.String()
    return scanner.token, scanner.error
}


func (scanner *Scanner) NextRune() rune {
    var r rune = 0
    r, _, scanner.error = scanner.reader.ReadRune()
    return r
}
