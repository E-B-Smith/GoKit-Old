

//----------------------------------------------------------------------------------------
//
//                                                                    ConfigurationUtil.go
//                                                  ServerUtil: Basic API server utilities
//
//                                                                   E.B.Smith, March 2016
//                        -©- Copyright © 2015-2016 Edward Smith, all rights reserved. -©-
//
//----------------------------------------------------------------------------------------


package ServerUtil


import (
    "os"
    "fmt"
    "net"
    "html"
    "path"
    "time"
    "bytes"
    "errors"
    "runtime"
    "strings"
    "syscall"
    "os/signal"
    "html/template"
    "runtime/pprof"
    "violent.blue/GoKit/Util"
    "violent.blue/GoKit/pgsql"
    "violent.blue/GoKit/Scanner"
    "violent.blue/GoKit/Log"
)


//----------------------------------------------------------------------------------------
//                                                              Open / Close Configuration
//----------------------------------------------------------------------------------------


func (config *Configuration) OpenConfig() error {
    Log.LogFunctionName()
    var error error

    //  Set up logging --

    Log.LogLevel = config.LogLevel
    Log.SetFilename(config.LogFilename)
    Log.LogTeeStderr = config.LogTeeStderr

    Log.Startf("%s version %s pid %d compiled %s.",
        config.ServiceName,
        CompileVersion(),
        os.Getpid(),
        CompileTime(),
    )
    Log.Debugf("Configuration: %+v.", config)

    //  Set our pid file --

    if error = config.CreatePIDFile(); error != nil {
        return error
    }

    //  Set our path --

    if error = os.Chdir(config.ServiceFilePath); error != nil {
        Log.Errorf("Error setting the home path '%s': %v.", config.ServiceFilePath, error)
        return error
    } else {
        config.ServiceFilePath, _ = os.Getwd()
        Log.Debugf("Working directory: '%s'", config.ServiceFilePath)
    }

    //  Load localized strings --

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

    //  Open the database --

    if error = config.ConnectDatabase(); error != nil {
        return error
    }

    return nil
}



func (config *Configuration) CloseConfig() {
    Log.LogFunctionName()
    config.DisconnectDatabase()
    config.DetachFromInterrupts()
    config.RemovePIDFile()
}


//  For use in template files
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
    result := fmt.Sprintf("%s PID %d Elapsed %s CPU %1.1f%% Thr %d/%d Mem %s Messages: %s",
        config.ServiceName,
        pinfo.PID,
        Util.HumanDuration(time.Since(pinfo.StartTime)),
        pinfo.CPUPercent,
        runtime.NumGoroutine(),
        runtime.NumCPU(),
        Util.HumanBytes(int64(pinfo.VMemory)),
        Util.HumanInt(int64(config.MessageCount)),
    )
    return result
}


func GetProfileNames() string {
    s := "Profiles:\n"
    profiles := pprof.Profiles()
    for _, profile := range profiles {
        s += profile.Name() + "\n"
    }
    return s
}


func GetGoprocs() string {
    buffer := new(bytes.Buffer)
    profile := pprof.Lookup("goroutine")
    profile.WriteTo(buffer, 1)
    s := "GoProcs:\n" + string(buffer.Bytes()) + "\n"
    return s
}


//  Process commands from the TCP pipe:  status | stop | help | hello | version
func ProcessTCPCommands(config *Configuration, connection net.Conn) {
    Log.LogFunctionName()
    defer connection.Close()
    Log.Infof("Accepted C&C connection from %s.", connection.RemoteAddr().String())
    helpString := ">>> Commands: 'status', 'stop', 'restart', 'help', 'version', 'profiles', 'goprocs'.\n"
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
                    s := fmt.Sprintf(">>> Software version %s.\n", CompileVersion())
                    _, error = connection.Write([]byte(s))
                case "status":
                    s := fmt.Sprintf("%s.\n", config.ServerStatusString())
                    _, error = connection.Write([]byte(s))
                case "profiles":
                    s := GetProfileNames()
                    _, error = connection.Write([]byte(s))
                case "goprocs":
                    s := GetGoprocs()
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
                go ProcessTCPCommands(config, connection)
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
    signal.Notify(
        config.signalChannel,
        syscall.SIGHUP,
        syscall.SIGINT,
        syscall.SIGKILL,
        syscall.SIGUSR1,
        syscall.SIGTERM,
    )
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
    Log.Infof("Starting database %s.", config.DatabaseURL)
    var error error
    config.PGSQL, error = pgsql.ConnectDatabase(config.DatabaseURL)
    if error != nil {
        Log.Errorf("Can't open database '%s':\n%v.", config.DatabaseURL, error)
        return error
    }
    pgsql.EnableInfiniteTime()
    config.DB = config.PGSQL.DB
    return nil
}


func (config *Configuration) DisconnectDatabase() {
    if config.PGSQL != nil {
        Log.Debugf("Stopping database %s.", config.DatabaseURL)
        config.PGSQL.DisconnectDatabase()
    }
    config.PGSQL = nil
    config.DB = nil
}


func (config *Configuration) DatabaseIsConnected() bool {
    return (config.PGSQL != nil)
}

