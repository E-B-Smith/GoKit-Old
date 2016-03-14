//  Configuration  -  Parse the configuration file.
//
//  E.B.Smith  -  March, 2015


package ServerUtil


import (
    "io"
    "os"
    "fmt"
    "net"
    "time"
    "html"
    "path"
    "errors"
    "strings"
    "syscall"
    "os/signal"
    "database/sql"
    "html/template"
    "violent.blue/GoKit/Log"
    "violent.blue/GoKit/Util"
    "violent.blue/GoKit/pgsql"
    "violent.blue/GoKit/Scanner"
)


type Configuration struct {
    SoftwareVersion string
    ServiceName     string
    ServicePort     int
    ServiceFilePath string
    ServicePrefix   string
    LogLevel        Log.LogLevelType
    LogFilename     string
    ServerURL       string
    WebLog          string

    TestingEnabled  bool

    AppName                 string
    AppLinkURL              string
    AppLinkScheme           string
    AppStoreURL             string
    ShortLinkURL            string
    LocalizationFile        string
    TemplatesPath           string

    Template                *template.Template

    //  Email

    EmailAddress        string  //  "beinghappy@beinghappy.io"
    EmailAccount        string  //  "beinghappy@beinghappy.io"
    EmailPassword       string  //  "*****"
    EmailSMTPServer     string  //  "smtp.gmail.com:587"

    //  Database

    DatabaseURI     string
    PGSQL           *pgsql.PGSQL
    DB              *sql.DB

    //  HappyPulse config

    PulsesAreFree       bool

    //  Global info

    MessageCount    int
    signalChannel   chan os.Signal
}


/*
func (config *Configuration) ServiceURL() string {
    if config.ServerPort == 80 || config.ServerPort == 0 {
        return config.HostURL + config.ServicePrefix
    } else {
        return fmt.Sprintf("%s:%d%s", config.HostURL, config.ServerPort, config.ServicePrefix)
    }
}
*/


func (config *Configuration) ServiceURL() string {
    return config.ServerURL + config.ServicePrefix
}



//----------------------------------------------------------------------------------------
//                                                                      ParseConfiguration
//----------------------------------------------------------------------------------------


func (configuration *Configuration) ParseFile(inputFile *os.File) error {
    var error error
    scanner := Scanner.NewFileScanner(inputFile)
    for !scanner.IsAtEnd() {
        //Log.Debugf("Token: '%s'.", scanner.Token())

        var identifier string
        identifier, error = scanner.ScanIdentifier()
        //Log.Debugf("Scanned '%s'.", scanner.Token())

        if error == io.EOF {
            return nil
        }
        if error != nil {
            return error
        }
        if identifier == "service-name" {
            configuration.ServiceName, error = scanner.ScanNext()
            if error != nil { return error }
            continue
        }
        if identifier == "service-port" {
            configuration.ServicePort, error = scanner.ScanInt()
            if error != nil { return error }
            continue
        }
        if identifier == "service-file-path" {
            configuration.ServiceFilePath, error = scanner.ScanNext()
            if error != nil { return error }
            continue
        }
        if identifier == "service-prefix" {
            configuration.ServicePrefix, error = scanner.ScanNext()
            if error != nil { return error }
            continue
        }
        if identifier == "database-uri" {
            configuration.DatabaseURI, error = scanner.ScanNext()
            if error != nil { return error }
            continue
        }
        if identifier == "web-log" {
            configuration.WebLog, error = scanner.ScanNext()
            if error != nil { return error }
            continue
        }
        if identifier == "log-level" {
            s, error := scanner.ScanNext()
            if error != nil { return error }
            if strings.HasPrefix(s, "Log") { s = s[3:] }
            configuration.LogLevel = Log.LogLevelFromString(s)
            if configuration.LogLevel == Log.LevelInvalid {
                return scanner.SetErrorMessage("Invalid log level")
            }
            continue
        }
        if identifier == "log-name" {
            configuration.LogFilename, error = scanner.ScanNext()
            if error != nil { return error }
            continue
        }
        if identifier == "server-url" {
            configuration.ServerURL, error = scanner.ScanNext()
            if error != nil { return error }
            continue
        }
        if identifier == "app-link-url" {
            configuration.AppLinkURL, error = scanner.ScanString()
            if error != nil { return error }
            continue
        }
        if identifier == "app-link-scheme" {
            configuration.AppLinkScheme, error = scanner.ScanString()
            if error != nil { return error }
            continue
        }
        if identifier == "app-name" {
            configuration.AppName, error = scanner.ScanString()
            if error != nil { return error }
            continue
        }
        if identifier == "short-link-url" {
            configuration.ShortLinkURL, error = scanner.ScanString()
            if error != nil { return error }
            continue
        }
        if identifier == "localization-file" {
            configuration.LocalizationFile, error = scanner.ScanString()
            if error != nil { return error }
            continue
        }
        if identifier == "app-store-url" {
            configuration.AppStoreURL, error = scanner.ScanString()
            if error != nil { return error }
            continue
        }
        if identifier == "templates-path" {
            configuration.TemplatesPath, error = scanner.ScanString()
            if error != nil { return error }
            continue
        }
        if identifier == "email-account" {
            configuration.EmailAccount, error = scanner.ScanString()
            if error != nil { return error }
            continue
        }
        if identifier == "email-password" {
            configuration.EmailPassword, error = scanner.ScanString()
            if error != nil { return error }
            continue
        }
        if identifier == "email-address" {
            configuration.EmailAddress, error = scanner.ScanQuotedString()
            if error != nil { return error }
            continue
        }
        if identifier == "email-smtp-server" {
            configuration.EmailSMTPServer, error = scanner.ScanString()
            if error != nil { return error }
            continue
        }
        if identifier == "pulses-are-free" {
            configuration.PulsesAreFree, error = scanner.ScanBool()
            if error != nil { return error }
            continue
        }
        if identifier == "testing-enabled" {
            configuration.TestingEnabled, error = scanner.ScanBool()
            if error != nil { return error }
            continue
        }
        if identifier == "app-link-redirect-url"    ||
           identifier == "client-app-min-version"   ||
           identifier == "client-app-min-data-date" ||
           identifier == "debug-mode-on"            ||
           identifier == "production-mode-on"{
            Log.Warningf("'%s' is deprecated.", identifier)
            if identifier == "client-app-min-data-date" {
                scanner.ScanTimestamp()
            } else {
                scanner.ScanString()
            }
            continue
        }
        return scanner.SetErrorMessage("Configuration identifier expected")
    }

    //  Check the config for basic correctness --

    if configuration.ServicePort == 0 { configuration.ServicePort = 80 }

    if  len(configuration.ServiceName) == 0     ||
        len(configuration.ServiceFilePath) == 0 ||
        len(configuration.ServicePrefix) == 0   ||
        len(configuration.DatabaseURI) == 0     ||
        len(configuration.ServerURL) == 0 {
        return errors.New("Missing config parameters")
    }

    //  Done --

    return error
}


func (config *Configuration) ParseFilename(filename string) error {
    inputFile, error := os.Open(filename)
    if error != nil {
        return fmt.Errorf("Error: Can't open file '%s' for reading: %v.", filename, error)
    }
    defer inputFile.Close()
    error = config.ParseFile(inputFile)
    if error != nil { return error }
    //Log.Debugf("Parsed configuration: %v.", config)
    return nil
}


func UnescapeString(args ...interface{}) string {
    Log.Debugf("UnescapeString:")
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


func (config *Configuration) ApplyConfiguration() error {

    //  Load localized strings --

    var error error = nil
    if len(config.LocalizationFile) > 0 {

        Log.Infof("Loading localized strings from %s.", config.LocalizationFile)

        error = config.LoadLocalizedStrings()
        if error != nil {
            error = fmt.Errorf("Can't open localization file: %v.", error)
            return error
        }
    }

    //  Load templates --

    if len(config.TemplatesPath) > 0 {

        Log.Infof("Loading templates from %s.", config.TemplatesPath)

        path := config.TemplatesPath+"/*"

        config.Template = template.New("Base")
        config.Template = config.Template.Funcs(template.FuncMap{"unescapeString": UnescapeString})
        config.Template, error = config.Template.ParseGlob(path)
        if error != nil || config.Template == nil {
            if error == nil { error = fmt.Errorf("No files.") }
            error = fmt.Errorf("Can't parse template files: %v.", error)
            return error
        }
    }

    return nil
}


//----------------------------------------------------------------------------------------
//                                                              PID File Service Functions
//----------------------------------------------------------------------------------------


func (config *Configuration) PIDFileName() string {
    name := "~/.run/" + config.ServiceName + ".pid"
    name = Util.AbsolutePath(name)
    return name
}


func (config *Configuration) CreatePIDFile() error {
    //  Try to create the pid file.
    //  -- On success, write pid to file.
    //  -- On Failure, see if pid is still running.
    //     If running, fail, else remove file and try again once.

    var dirPerm  os.FileMode = 0750
    var filePerm os.FileMode = 0640

    filename := config.PIDFileName()
    Log.Debugf("PID file: %s permissions: %d %o %x.", filename, dirPerm, dirPerm, dirPerm)
    error := os.MkdirAll(path.Dir(filename), dirPerm)
    if error != nil {
        Log.Warningf("Can't create pid directory %s: %v.", filename, error)
        return error
    }

    // actualPerm, _ := os.Stat(path.Dir(filename))
    // Log.Debugf("Dir: %s Perm: %o.", actualPerm.Name(), actualPerm.Mode())

    file, error := os.OpenFile(filename, os.O_WRONLY | os.O_CREATE | os.O_EXCL, filePerm)
    if error == nil {
        pinfo, error := Util.GetProcessInfo(os.Getpid())
        if error == nil {
            fmt.Fprintf(file, "%d\t%s\n", pinfo.PID, pinfo.Command)
            file.Close()
            return nil
        }
    }
    Log.Debugf("Can't create PID file: %v.", error)
    file, error = os.OpenFile(filename, os.O_RDONLY, filePerm)
    defer file.Close()
    if error != nil { return error }

    scanner := Scanner.NewFileScanner(file)
    pid, _ := scanner.ScanInt()
    command, _ := scanner.ScanToEOL()
    Log.Debugf("PID file contents: %d %s.", pid, command)

    pinfo, error := Util.GetProcessInfo(pid)
    if error != nil || pinfo.Command != command {
        Log.Warningf("Removing old pid file...")
        os.Remove(filename)
        return config.CreatePIDFile()
    }

    return errors.New("Already running")
}


func (config *Configuration) RemovePIDFile() error {
    //  Remove the pid file.
    filename := config.PIDFileName()
    Log.Debugf("Removing PID file %s.", filename)
    return os.Remove(filename)
}



//----------------------------------------------------------------------------------------
//                                                           TCP Command & Control Channel
//----------------------------------------------------------------------------------------


func (config *Configuration) ServerStatusString() string {
    pinfo, _ := Util.GetProcessInfo(os.Getpid())
    result := fmt.Sprintf("%s PID %d Elapsed %s CPU %1.1f%% Mem %s Messages: %s",
        config.ServiceName,
        pinfo.PID,
        Util.HumanDuration(time.Since(pinfo.StartTime)),
        pinfo.CPUPercent,
        Util.HumanBytes(int64(pinfo.VMemory)),
        Util.HumanInt(int64(config.MessageCount)),
    )
    return result
}


func ProcessCommands(config *Configuration, connection net.Conn) {
    //  Commands:  status | stop | help | hello | version
    Log.LogFunctionName()
    defer connection.Close()
    Log.Infof("Accepted C&C connection from %s.", connection.RemoteAddr().String())
    helpString := ">>> Commands: 'status', 'stop', 'restart', 'help', 'version'.\n"
    var error error
    timeout, error := time.ParseDuration("30s")
    if error != nil { Log.Errorf("Error parsing duration: %v.", error) }

    var commandBuffer string
    buffer := make([]byte, 256)
    for error == nil {
        connection.SetDeadline(time.Now().Add(timeout))
        var n int
        n, error = connection.Read(buffer)
        if n <= 0 && error != nil { break }
        Log.Debugf("Read %d characters.", n)

        commandBuffer += string(buffer[:n])
        index := strings.Index(commandBuffer, "\n")
        for index > -1 && error == nil {
            command := strings.ToLower(commandBuffer[:index])
            command  = strings.TrimSpace(command)
            if index < len(commandBuffer)-1 {
                commandBuffer = commandBuffer[index+1:]
            } else {
                commandBuffer = ""
            }
            index = strings.Index(commandBuffer, "\n")

            Log.Infof("C&C command '%s'.", command)
            switch command {
                case "hello":
                    _, error = connection.Write([]byte(">>> Hello.\n"))
                case "version":
                    s := fmt.Sprintf(">>> Software version %s.\n", config.SoftwareVersion)
                    _, error = connection.Write([]byte(s))
                case "status":
                    s := fmt.Sprintf("%s.\n", config.ServerStatusString())
                    _, error = connection.Write([]byte(s))
                case "stop":
                    _, error = connection.Write([]byte(">>> Stopping.\n"))
                    myself, _ := os.FindProcess(os.Getpid())
                    myself.Signal(syscall.SIGHUP)
                case "", " ", "\n":
                case "help", "?", "h":
                    _, error = connection.Write([]byte(helpString))
                default:
                    message := fmt.Sprintf(">>> Unknown command '%s'.\n", command)
                    _, error = connection.Write([]byte(message))
                    if error != nil { break; }
                    _, error = connection.Write([]byte(helpString))
            }
        }
    }
    if error != nil {
        Log.Debugf("Connection closed with error %v.", error)
    } else {
        Log.Debugf("Connection closed without error.")
    }
}


func (config *Configuration) StartTCPCommandChannel() {
    //  Start listening for commands --
    Log.LogFunctionName()
    port := fmt.Sprintf("localhost:%d", config.ServicePort+1)
    listener, error := net.Listen("tcp", port)
    if error != nil {
        Log.LogError(error)
        return
    }
    go func() {
        defer listener.Close()
        Log.Infof("Listening for C&C connections on %s.", listener.Addr().String())
        for {
            connection, error := listener.Accept()
            if error != nil {
                Log.LogError(error)
                break
            } else {
                go ProcessCommands(config, connection)
            }
        }
    } ()
}



//----------------------------------------------------------------------------------------
//                                                             Interrupt Handler Functions
//----------------------------------------------------------------------------------------


func (config *Configuration) AttachToInterrupts(httpListener net.Listener) {
    //  Set up an interrupt handler --
    config.signalChannel = make(chan os.Signal, 1)
    signal.Notify(config.signalChannel,
        syscall.SIGHUP, syscall.SIGINT, syscall.SIGKILL, syscall.SIGUSR1, syscall.SIGTERM)
    go func() {
        for signal := range config.signalChannel {
            fmt.Fprintf(os.Stderr, "Signal %v\n", signal)
            if signal == syscall.SIGUSR1 {
                statusString := config.ServerStatusString()
                Log.Infof("%s", statusString)
            } else {
                Log.Infof("Quit signal received.")
                httpListener.Close()
            }
        }
    } ()
    //  Now start our TCP command channel --
    config.StartTCPCommandChannel()
}


func (config *Configuration) DetachFromInterrupts() {
    if config.signalChannel != nil {
        signal.Stop(config.signalChannel)
    }
}



//----------------------------------------------------------------------------------------
//                                                                                Database
//----------------------------------------------------------------------------------------


func (config *Configuration) ConnectDatabase() error {
    Log.Infof("Starting database %s.", config.DatabaseURI)
    var error error
    config.PGSQL, error = pgsql.ConnectDatabase(config.DatabaseURI)
    if error != nil {
        Log.Errorf("Can't open database '%s':\n%v.", config.DatabaseURI, error)
        return error
    }
    pgsql.EnableInfiniteTime()
    config.DB = config.PGSQL.DB
    return nil
}


func (config *Configuration) DisconnectDatabase() {
    if config.PGSQL != nil {
        Log.Debugf("Stopping database %s.", config.DatabaseURI)
        config.PGSQL.DisconnectDatabase()
    }
    config.PGSQL = nil
    config.DB = nil
}


func (config *Configuration) DatabaseIsConnected() bool {
    return (config.PGSQL != nil)
}


//----------------------------------------------------------------------------------------
//                                                                                   Close
//----------------------------------------------------------------------------------------


func (config *Configuration) Close() {
    Log.Debugf("Cleaning up config.")
    config.DisconnectDatabase()
    config.RemovePIDFile()
}

