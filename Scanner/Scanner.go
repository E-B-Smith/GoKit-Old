//  Scanner  -  A parser to parse a configuration file.
//
//  E.B.Smith  -  November, 2014


package Scanner


import (
    "os"
    "fmt"
    "bytes"
    "bufio"
    "unicode"
    "errors"
    "strconv"
)


type Scanner struct {
    file        *os.File
    reader      *bufio.Reader
    lineNumber  int
    error       error
    token       string
}


func NewScanner(file *os.File) *Scanner {
    if file == nil { return nil }
    scanner := new(Scanner)
    scanner.file = file
    scanner.reader = bufio.NewReader(file)
    scanner.lineNumber = 1
    return scanner
}


func (scanner *Scanner) FileName() string {
    return scanner.file.Name()
}


func (scanner *Scanner) LineNumber() int {
    return scanner.lineNumber
}


func (scanner *Scanner) IsAtEnd() bool {
    return scanner.error != nil
}


func (scanner *Scanner) Token() string {
    return scanner.token;
}


//  Scan Routines --


func IsValidIdentifierStartRune(r rune) bool {
    return unicode.IsLetter(r) || r == '_'
    }


func IsValidIdentifierRune(r rune) bool {
    return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_'
    }

func IsOctalDigit(r rune) bool {
    return unicode.IsDigit(r) && r != '8' && r != '9'
    }

func ZIsSpace(r rune) bool {
    return unicode.IsSpace(r) || r == '#'
    }


func ZIsLineFeed(r rune) bool {
    return r == '\n' || r == '\u0085'
    }


func (scanner *Scanner) ScanSpaces() (token string, error error) {
    for ! scanner.IsAtEnd() {
        var r rune
        r, _, scanner.error = scanner.reader.ReadRune()

        if r == '#' {
            for !scanner.IsAtEnd() && !ZIsLineFeed(r) {
                r, _, scanner.error = scanner.reader.ReadRune()
                }
            }

        if ZIsLineFeed(r) {
            scanner.lineNumber++
            continue
            }

        if ZIsSpace(r) {
            continue
            }

        scanner.reader.UnreadRune()
        return "", nil
        }

    return "", scanner.error;
    }


func IsValidStringRune(r rune) bool {
    if r == ';' || r == ',' || ZIsSpace(r) { return false }
    return unicode.IsGraphic(r)
    }


func (scanner *Scanner) ScanString() (next string, error error) {
    _, error = scanner.ScanSpaces()

    var (r rune; buffer bytes.Buffer)
    r, _, scanner.error = scanner.reader.ReadRune()

    for IsValidStringRune(r) {
        buffer.WriteRune(r)
        r, _, scanner.error = scanner.reader.ReadRune()
        }
    scanner.reader.UnreadRune()

    scanner.token = buffer.String()
    return scanner.token, nil
    }


func (scanner *Scanner) ScanInteger() (int int, error error) {
    scanner.ScanSpaces()
    var r rune
    r, _, scanner.error = scanner.reader.ReadRune()

    if ! unicode.IsDigit(r) {
        scanner.reader.UnreadRune()
        scanner.token, _ = scanner.ScanNext()
        return 0, errors.New("Integer expected")
        }

    var buffer bytes.Buffer
    for unicode.IsDigit(r) {
        buffer.WriteRune(r)
        r, _, scanner.error = scanner.reader.ReadRune()
        }
    scanner.reader.UnreadRune()

    scanner.token = buffer.String()
    return  strconv.Atoi(scanner.token)
    }


func (scanner *Scanner) ScanNext() (next string, error error) {
    scanner.ScanSpaces()
    var r rune
    r, _, scanner.error = scanner.reader.ReadRune()

    if r == '"' {
        return scanner.ScanQuotedString()
        }

    if unicode.IsPunct(r) {
        var buffer bytes.Buffer
        buffer.WriteRune(r)
        scanner.token = buffer.String()
        return scanner.token, nil
        }

    if unicode.IsDigit(r) {
        scanner.reader.UnreadRune()
        scanner.ScanInteger()
        return scanner.token, scanner.error
        }

    scanner.reader.UnreadRune()
    return scanner.ScanString()
    }


func (scanner *Scanner) ScanIdentifier() (identifier string, error error) {
    scanner.ScanSpaces()
    var r rune
    r, _, scanner.error = scanner.reader.ReadRune()
    if scanner.error != nil {
        return "", scanner.error
        }

    if ! IsValidIdentifierStartRune(r) {
        scanner.reader.UnreadRune()
        scanner.ScanNext()
        return "", errors.New("Identifier expected")
        }

    var buffer bytes.Buffer
    for IsValidIdentifierRune(r) {
        buffer.WriteRune(r)
        r, _, scanner.error = scanner.reader.ReadRune()
        }
    scanner.reader.UnreadRune()

    scanner.token = buffer.String()
    return scanner.token, nil
    }


func (scanner *Scanner) ScanOctal() (Integer int, error error) {
    scanner.ScanSpaces()
    var r rune
    r, _, scanner.error = scanner.reader.ReadRune()

    if ! IsOctalDigit(r) {
        scanner.reader.UnreadRune()
        scanner.token, _ = scanner.ScanNext()
        return 0, errors.New("Octal number expected")
        }

    var buffer bytes.Buffer
    for IsOctalDigit(r) {
        buffer.WriteRune(r)
        r, _, scanner.error = scanner.reader.ReadRune()
        }
    scanner.reader.UnreadRune()

    scanner.token = buffer.String()
    val, error := strconv.ParseInt(scanner.token, 8, 0)
    return int(val), error
}


func (scanner *Scanner) ScanQuotedString() (string, error) {
    scanner.ScanSpaces()
    var r rune
    r, _, scanner.error = scanner.reader.ReadRune()

    if r != '"' {
        scanner.token = string(r)
        return "", errors.New("Quoted string expected")
        }
    scanner.reader.UnreadRune()

    var result string
    parseCount, error := fmt.Fscanf(scanner.reader, "%q", result)
    if error != nil { return "", error }
    if parseCount != 1 { return "", errors.New("Quoted string expected") }

    return result, nil
}

