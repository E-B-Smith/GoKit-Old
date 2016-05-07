

//----------------------------------------------------------------------------------------
//
//                                                                        Configuration.go
//                                                  ServerUtil: Basic API server utilities
//
//                                                                   E.B.Smith, March 2016
//                        -©- Copyright © 2015-2016 Edward Smith, all rights reserved. -©-
//
//----------------------------------------------------------------------------------------


package ServerUtil


import (
    "io"
    "os"
    "fmt"
    "reflect"
    "strings"
    "unicode"
    "database/sql"
    "html/template"
    "violent.blue/GoKit/pgsql"
    "violent.blue/GoKit/Scanner"
    "violent.blue/GoKit/Log"
)


type Configuration struct {

    ServiceName     string
    ServicePort     int
    ServiceFilePath string
    ServicePrefix   string
    ServerURL       string

    TestingEnabled  bool

    //  Logging --

    LogLevel        Log.LogLevelType `enum:"LogLevelInvalid,LogLevelAll,LogLevelDebug,LogLevelStart,LogLevelExit,LogLevelInfo,LogLevelWarning,LogLevelError"`
    LogTeeStderr    bool
    LogFilename     string

    //  Database --

    DatabaseURL     string
    PGSQL           *pgsql.PGSQL
    DB              *sql.DB

    //  For app deep links --

    AppName                 string
    AppLinkURL              string
    AppLinkScheme           string
    AppStoreURL             string
    ShortLinkURL            string

    //  Locaization and templates --

    LocalizationFile        string
    TemplatesPath           string

    Template                *template.Template

    //  Email

    EmailAddress        string  //  "beinghappy@beinghappy.io"
    EmailAccount        string  //  "beinghappy@beinghappy.io"
    EmailPassword       string  //  "*****"
    EmailSMTPServer     string  //  "smtp.gmail.com:587"

    //  Global server stats / info --

    MessageCount    int
    signalChannel   chan os.Signal
}


//----------------------------------------------------------------------------------------
//                                                                     CompileTime/Version
//----------------------------------------------------------------------------------------


var compileVersion              string = "0.0.0"
var compileTime                 string = "Sun Mar 6 09:01:25 PST 2016"

func CompileVersion() string    { return compileVersion }
func CompileTime() string       { return compileTime }

func (config *Configuration) ServiceURL() string {
    return config.ServerURL + config.ServicePrefix
}


//----------------------------------------------------------------------------------------
//                                                                      ParseConfiguration
//----------------------------------------------------------------------------------------


func (config *Configuration) ParseConfigFileNamed(filename string) error {
    inputFile, error := os.Open(filename)
    if error != nil {
        return fmt.Errorf("Error: Can't open file '%s' for reading: %v.", filename, error)
    }
    defer inputFile.Close()
    error = config.ParseConfigFile(inputFile)
    if error != nil { return error }
    //Log.Debugf("Parsed configuration: %v.", config)
    return nil
}


func (config *Configuration) ParseConfigFile(inputFile *os.File) error {

    //  Set some default values --

    config.ServicePort = 80

    //  Start relecting --

    configStructPtrValue := reflect.ValueOf(config)
    Log.Debugf("Kind: %+v Value: %+v", configStructPtrValue.Kind(), configStructPtrValue)
    configStructValue := configStructPtrValue.Elem()
    Log.Debugf("Kind: %+v Value: %+v", configStructValue.Kind(), configStructValue)

    scanner := Scanner.NewFileScanner(inputFile)
    for !scanner.IsAtEnd() {
        var error error

        var identifier string
        identifier, error = scanner.ScanIdentifier()
        Log.Debugf("Scanned '%s'.", scanner.Token())

        if error == io.EOF {
            break
        }
        if error != nil {
            return error
        }

        //  Find the identifier --

        fieldName := CamelCaseFromIdentifier(identifier)
        Log.Debugf("FieldName: '%s'.", fieldName)
        field := configStructValue.FieldByName(fieldName)
        if ! field.IsValid() {
            return scanner.SetErrorMessage("Configuration identifier expected")
        }
        structField, _ := configStructValue.Type().FieldByName(fieldName)

        var (
            i int64
            s string
            b bool
            f float64
        )

        switch field.Type().Kind() {

        case reflect.Bool:
            b, error = scanner.ScanBool()
            if error != nil { return error }
            field.SetBool(b)

        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

            enumValues := structField.Tag.Get("enum")
            if len(enumValues) > 0 {

                s, error = scanner.ScanNext()
                if error != nil {
                    scanner.SetError(error)
                    return scanner.LastError()
                }

                i, error = enumFromString(s, enumValues)
                if error != nil { return error }

            } else {

                i, error = scanner.ScanInt64()
                if error != nil { return error }

            }
            field.SetInt(i)

        case reflect.Float32, reflect.Float64:
            f, error = scanner.ScanFloat64()
            if error != nil { return error }
            field.SetFloat(f)

        case reflect.String:
            s, error = scanner.ScanNext()
            if error != nil { return error }
            field.SetString(s)

        default:
            return fmt.Errorf("Error: '%s' unhandled type: %s", identifier, field.Type().Name())
        }
    }

    //  Check for basic values --

    var err error

    checkNotZero := func(name string) {
        fieldValue := config.ValueByName(name)
        if fieldValue.IsValid() {
            if fieldValue.Type().Kind() == reflect.String {
                if fieldValue.String() != "" { return }
            } else {
                if fieldValue.Int() != 0 { return }
            }
        }
        if err != nil { err = fmt.Errorf("Missing config parameter: %s", name) }
    }


    config.MessageCount = 0
    checkNotZero("ServiceName")
    checkNotZero("ServiceFilePath")
    checkNotZero("ServicePrefix")
    checkNotZero("DatabaseURI")
    checkNotZero("ServerURL")

    //  Done --

    return err
}


func (config *Configuration) ValueByName(name string) reflect.Value {
    return reflect.ValueOf(config).Elem().FieldByName(name)
}


func enumFromString(s string, enumValues string) (int64, error) {

    enumArray := make([]string, 0)
    a := strings.Split(enumValues, ",")
    for _, enum := range a {
        enum = strings.TrimSpace(enum)
        if len(enum) > 0 {
            enumArray = append(enumArray, enum)
        }
    }

    for i, val := range enumArray {
        if val == s { return int64(i), nil }
    }

    return -1, fmt.Errorf("Invalid enum '%s'", s)
}


//----------------------------------------------------------------------------------------
//                                                                 CamelCaseFromIdentifier
//----------------------------------------------------------------------------------------


//  Return possible Golang camel-case variations on the string.
func CamelCaseFromIdentifier(s string) string {

    lastWasUpper := false
    words := make([]string, 0, 5)
    var word []rune
    for _, r := range s {

        switch {

        case r == '-' || r == '_':
            words = append(words, string(word))
            word = word[:0]
            lastWasUpper = false

        case unicode.IsUpper(r):
            if ! lastWasUpper {
                words = append(words, string(word))
                word = word[:0]
            }
            word = append(word, r)
            lastWasUpper = true

        default:
            word = append(word, r)
            lastWasUpper = false
        }
    }
    if len(word) > 0 {
        words = append(words, string(word))
    }


    //  String together the parts.  Upper-case any special words:

    upperWords := map[string]bool {
        "http": true,
        "https":true,
        "url":  true,
        "uri":  true,
        "urn":  true,
        "smtp": true,
        "xml":  true,
        "json": true,
        "id":   true,
    }

    var camelString string
    for _, part := range words {

        part = strings.ToLower(part)
        if _, ok := upperWords[part]; ok {
            part = strings.ToUpper(part)
        } else {
            part = strings.Title(part)
        }


        camelString += part
    }

    return camelString
}


