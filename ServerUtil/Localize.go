//  Configuration  -  Parse the configuration file.
//
//  E.B.Smith  -  March, 2015


package ServerUtil


import (
    "os"
    "io"
    "fmt"
    "violent.blue/GoKit/Log"
    "violent.blue/GoKit/Scanner"
)


//----------------------------------------------------------------------------------------
//
//                                                                       Localized Strings
//
//                                                                               Localizef
//
//----------------------------------------------------------------------------------------


var globalLocalizeStringMap = make(map[string]*string)


func (config *Configuration) Localizef(messageKey, defaultFormat string, args ...interface{}) string {

    //  Check the string map first.  If not found use the default.

    var (format *string; ok bool)
    if format, ok = globalLocalizeStringMap[messageKey]; ! ok {
        format = &defaultFormat
    }
    return fmt.Sprintf(*format, args...)
}



//----------------------------------------------------------------------------------------
//
//                                                                       Localized Strings
//
//                                                                               Localizef
//
//----------------------------------------------------------------------------------------


func (config *Configuration) LoadLocalizedStrings() error {
    Log.LogFunctionName()

    if len(config.LocalizationFile) <= 0 {
        Log.Warningf("No filename for localized string file set.")
    }

    localizedMap := make(map[string]*string)

    var (
        error error
        file *os.File
        filePerm os.FileMode = 0640
    )
    file, error = os.OpenFile(config.LocalizationFile, os.O_RDONLY, filePerm)
    if error != nil {
        Log.Errorf("Can't open localized string file: %v.", error)
        return error
    }
    scanner := Scanner.NewFileScanner(file)
    for !scanner.IsAtEnd() {

        var identifier string
        identifier, error = scanner.ScanIdentifier()
        if error != nil { break }

        var symbol string
        symbol, error = scanner.ScanNext()
        if error != nil || symbol != "=" {
            return scanner.SetErrorMessage("Equal sign expected")
        }

        var localString string
        localString, error = scanner.ScanQuotedString()
        if error != nil || symbol != "=" {
            return scanner.SetErrorMessage("Localized string expected")
        }

        symbol, error = scanner.ScanNext()
        if error != nil || symbol != ";" {
            return scanner.SetErrorMessage("Semi-colon expected")
        }

        localizedMap[identifier] = &localString
    }

    if error == io.EOF { error = nil }
    globalLocalizeStringMap = localizedMap

    return error
}

