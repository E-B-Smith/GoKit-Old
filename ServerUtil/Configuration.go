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
    "path"
    "errors"
    "strings"
    "syscall"
    "os/signal"
    "database/sql"
    "../log"
    "../pgsql"
    "../Scanner"
)


type Configuration struct {
    SoftwareVersion string
    ServiceName     string
    ServicePort     int
    ServiceFilePath string
    ServicePrefix   string
    IsInDebugMode   bool
    IsInProductionMode bool
    LogLevel        log.LogLevelType
    LogFilename     string
    ServerURL       string
    WebLog          string
    AppLinkRedirectURL  string

    //  Database

    DatabaseURI     string
    PGSQL           *pgsql.PGSQL
    DB              *sql.DB

    //  Global info

    MessageCount    int
    signalChannel   chan os.Signal

    //  Client configuration --

    ClientAppMinDataDate    time.Time;
    ClientAppMinVersion     string;
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
        //log.Debugf("Token: '%s'.", scanner.Token())

        var identifier string
        identifier, error = scanner.ScanIdentifier()
        //log.Debugf("Scanned '%s'.", scanner.Token())

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
        if identifier == "debug-mode-on" {
            configuration.IsInDebugMode, error = scanner.ScanBool()
            if error != nil { return error }
            continue
        }
        if identifier == "production-mode-on" {
            configuration.IsInProductionMode, error = scanner.ScanBool()
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
            configuration.LogLevel = log.LogLevelFromString(s)
            if configuration.LogLevel == log.LevelInvalid {
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
        if identifier == "client-app-min-data-date" {
            configuration.ClientAppMinDataDate, error = scanner.ScanTimestamp()
            if error != nil { return error }
            continue
        }
        if identifier == "client-app-min-version" {
            configuration.ClientAppMinVersion, error = scanner.ScanString()
            if error != nil { return error }
            continue
        }
        if identifier == "app-link-redirect-url" {
            configuration.AppLinkRedirectURL, error = scanner.ScanString()
            if error != nil { return error }
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

    return nil
}


func (config *Configuration) ParseFilename(filename string) error {
    inputFile, error := os.Open(filename)
    if error != nil {
        return fmt.Errorf("Error: Can't open file '%s' for reading: %v.", filename, error)
    }
    defer inputFile.Close()
    error = config.ParseFile(inputFile)
    if error != nil { return error }
    log.Debugf("Parsed configuration: %v.", config)
    return nil
}



//----------------------------------------------------------------------------------------
//                                                              PID File Service Functions
//----------------------------------------------------------------------------------------


func (config *Configuration) PIDFileName() string {
    name := "~/.run/" + config.ServiceName + ".pid"
    name = CleanupPath(name)
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
    log.Debugf("PID file: %s permissions: %d %o %x.", filename, dirPerm, dirPerm, dirPerm)
    error := os.MkdirAll(path.Dir(filename), dirPerm)
    if error != nil {
        log.Warningf("Can't create pid directory %s: %v.", filename, error)
        return error
    }

    // actualPerm, _ := os.Stat(path.Dir(filename))
    // log.Debugf("Dir: %s Perm: %o.", actualPerm.Name(), actualPerm.Mode())

    file, error := os.OpenFile(filename, os.O_WRONLY | os.O_CREATE | os.O_EXCL, filePerm)
    if error == nil {
        pinfo, error := GetProcessInfo(os.Getpid())
        if error == nil {
            fmt.Fprintf(file, "%d\t%s\n", pinfo.PID, pinfo.Command)
            file.Close()
            return nil
        }
    }
    log.Debugf("Can't create PID file: %v.", error)
    file, error = os.OpenFile(filename, os.O_RDONLY, filePerm)
    defer file.Close()
    if error != nil { return error }

    scanner := Scanner.NewFileScanner(file)
    pid, _ := scanner.ScanInt()
    command, _ := scanner.ScanToEOL()
    log.Debugf("PID file contents: %d %s.", pid, command)

    pinfo, error := GetProcessInfo(pid)
    if error != nil || pinfo.Command != command {
        log.Warningf("Removing old pid file...")
        os.Remove(filename)
        return config.CreatePIDFile()
    }

    return errors.New("Already running")
}


func (config *Configuration) RemovePIDFile() error {
    //  Remove the pid file.
    filename := config.PIDFileName()
    log.Debugf("Removing PID file %s.", filename)
    return os.Remove(filename)
}



//----------------------------------------------------------------------------------------
//                                                           TCP Command & Control Channel
//----------------------------------------------------------------------------------------


func (config *Configuration) ServerStatusString() string {
    pinfo, _ := GetProcessInfo(os.Getpid())
    result := fmt.Sprintf("%s PID %d Elapsed %s CPU %1.1f%% Mem %s Messages: %s",
        config.ServiceName,
        pinfo.PID,
        time.Since(pinfo.StartTime).String(),
        pinfo.CPUPercent,
        HumanBytes(int64(pinfo.VMemory)),
        HumanInt(int64(config.MessageCount)),
    )
    return result
}


func ProcessCommands(config *Configuration, connection net.Conn) {
    //  Commands:  status | stop | help | hello | version
    log.LogFunctionName()
    defer connection.Close()
    log.Infof("Accepted C&C connection from %s.", connection.RemoteAddr().String())
    helpString := ">>> Commands: 'status', 'stop', 'restart', 'help', 'version'.\n"
    var error error
    timeout, error := time.ParseDuration("30s")
    if error != nil { log.Errorf("Error parsing duration: %v.", error) }

    var commandBuffer string
    buffer := make([]byte, 256)
    for error == nil {
        connection.SetDeadline(time.Now().Add(timeout))
        var n int
        n, error = connection.Read(buffer)
        if n <= 0 && error != nil { break }
        log.Debugf("Read %d characters.", n)

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

            log.Infof("C&C command '%s'.", command)
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
        log.Debugf("Connection closed with error %v.", error)
    } else {
        log.Debugf("Connection closed without error.")
    }
}


func (config *Configuration) StartTCPCommandChannel() {
    //  Start listening for commands --
    log.LogFunctionName()
    port := fmt.Sprintf("localhost:%d", config.ServicePort+1)
    listener, error := net.Listen("tcp", port)
    if error != nil {
        log.LogError(error)
        return
    }
    go func() {
        defer listener.Close()
        log.Infof("Listening for C&C connections on %s.", listener.Addr().String())
        for {
            connection, error := listener.Accept()
            if error != nil {
                log.LogError(error)
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
                log.Infof("%s", statusString)
            } else {
                log.Infof("Quit signal received.")
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
    log.Infof("Starting database %s.", config.DatabaseURI)
    var error error
    config.PGSQL, error = pgsql.ConnectDatabase(config.DatabaseURI)
    if error != nil {
        log.Errorf("Can't open database '%s':\n%v.", config.DatabaseURI, error)
        return error
    }
    pgsql.EnableInfiniteTime()
    config.DB = config.PGSQL.DB
    return nil
}


func (config *Configuration) DisconnectDatabase() {
    if config.PGSQL != nil {
        log.Debugf("Stopping database %s.", config.DatabaseURI)
        config.PGSQL.DisconnectDatabase()
    }
    config.PGSQL = nil
    config.DB = nil
}



//----------------------------------------------------------------------------------------
//                                                                                   Close
//----------------------------------------------------------------------------------------


func (config *Configuration) Close() {
    log.Debugf("Cleaning up config.")
    config.DisconnectDatabase()
    config.RemovePIDFile()
}


//----------------------------------------------------------------------------------------
//                                                                       Localized Strings
//
//                                                                               Localizef
//
//----------------------------------------------------------------------------------------


func (config *Configuration) Localizef(format string, args ...interface{}) string {
    return fmt.Sprintf(format, args...)
}


