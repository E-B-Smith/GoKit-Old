//  RunShellCommand  -  Run a shell command.
//
//  E.B.Smith  -  March, 2015


package Util


import (
    "io"
    "sync"
    "time"
    "bytes"
    "bufio"
    "errors"
    "strings"
    "strconv"
    "os/exec"
    "text/scanner"
    "violent.blue/GoKit/Log"
)



//----------------------------------------------------------------------------------------
//                                                                         RunShellCommand
//----------------------------------------------------------------------------------------


func RunShellCommand(shellCommand string, parameters []string, standardIn []byte)    (standardOut []byte, standardError []byte, error error) {
    //
    //  Run a shell command --
    //

    path, error := exec.LookPath(shellCommand)
    if error != nil {
        Log.Errorf("Can't find path for '%s': %v.", shellCommand, error)
        return nil, nil, error
    }
    Log.Debugf("     Path: %s.", path)
    Log.Debugf("StdIn Len: %d.", len(standardIn))

    command := exec.Command(path, parameters...)

    var errorpipe *io.PipeReader
    errorpipe, command.Stderr = io.Pipe()

    var waiter sync.WaitGroup
    waiter.Add(2)
    go func() {
        buffer := make([]byte, 512)
        reader := bufio.NewReader(errorpipe)
        count, error := reader.Read(buffer)
        //Log.Debugf("Read %d bytes.  Error: %v.", count, error)
        for error == nil {
            standardError = append(standardError, buffer[:count]...)
            count, error = reader.Read(buffer)
            //Log.Debugf("Read %d bytes.  Error: %v.", count, error)
        }
        waiter.Done()
    } ()

    var pipeout *io.PipeReader
    pipeout, command.Stdout = io.Pipe()
    go func () {
        buffer := make([]byte, 512)
        reader := bufio.NewReader(pipeout)
        count, error := reader.Read(buffer)
        //Log.Debugf("Read %d bytes.  Error: %v.", count, error)
        for  error == nil {
             standardOut = append(standardOut, buffer[:count]...)
             count, error = reader.Read(buffer)
            //Log.Debugf("Read %d bytes.  Error: %v.", count, error)
        }
        waiter.Done()
    } ()

    var pipein io.WriteCloser

    if  len(standardIn) > 0 {
        pipein, error = command.StdinPipe()
        if error != nil {
            Log.Errorf("Can't open input pipe: %v.", error)
            return nil, nil, error
        }

        _, error = pipein.Write([]byte(standardIn))
        if error != nil {
            Log.Errorf("Error writing stdin: %v.", error)
            return nil, nil, error
        }
    }

    error = command.Start()
    if error != nil {
        Log.Errorf("Error starting '%s': %v.", path, error)
        return nil, nil, error
    }

    if pipein != nil { pipein.Close() }
    error = command.Wait()
    errorpipe.Close()
    pipeout.Close()
    waiter.Wait()

    if error != nil {
        Log.Errorf("Script wait error: %v.", error)
    }

    return standardOut, standardError, error
}



//----------------------------------------------------------------------------------------
//                                                                          GetProcessInfo
//----------------------------------------------------------------------------------------


type ProcessInfo struct {
    PID         int         //  pid
    StartTime   time.Time   //  lstart
    CPUTime     float32     //  time
    CPUPercent  float32     //  pcpu
    MessagesIn  int         //  msgrcv
    MessagesOut int         //  msgsnd
    BlocksIn    int         //  inblk
    BlocksOut   int         //  oublk
    VMemory     int         //  vsize
    Command     string      //  comm
    CommandLine string      //  command
}


func GetProcessInfo(PID int) (*ProcessInfo, error) {
    pidStr := strconv.Itoa(PID)
    params := []string{ "-opid=,pcpu=,msgrcv=,msgsnd=,inblk=,oublk=,vsize=,time=", "-p", pidStr }

    stdout, stderr, error := RunShellCommand("ps", params, nil)
    if error != nil { return nil, error }
    if len(stderr) > 0 { return nil, errors.New(string(stderr)) }
    // fmt.Printf("%s\n", stdout)

    buffer := bytes.NewBufferString(string(stdout))
    var scan scanner.Scanner
    scan.Init(buffer)
    pinfo := new(ProcessInfo)

    //  Do the numbers --

    scan.Scan()
    pinfo.PID, _ = strconv.Atoi(scan.TokenText())

    scan.Scan()
    f, _ := strconv.ParseFloat(scan.TokenText(), 32)
    pinfo.CPUPercent = float32(f)

    scan.Scan()
    pinfo.MessagesIn, _ = strconv.Atoi(scan.TokenText())

    scan.Scan()
    pinfo.MessagesOut, _ = strconv.Atoi(scan.TokenText())

    scan.Scan()
    pinfo.BlocksIn, _ = strconv.Atoi(scan.TokenText())

    scan.Scan()
    pinfo.BlocksOut, _ = strconv.Atoi(scan.TokenText())

    scan.Scan()
    pinfo.VMemory, _ = strconv.Atoi(scan.TokenText())

    scan.Scan()
    min, _ := strconv.Atoi(scan.TokenText())

    scan.Scan()
    scan.Scan()
    sec, _ := strconv.ParseFloat(scan.TokenText(), 32)
    pinfo.CPUTime = float32(min) * 60.0 + float32(sec)

    //  Do start time --

    stdout, stderr, error = RunShellCommand("ps", []string{"-olstart=", "-p", pidStr}, nil)
    if error != nil { return nil, error }
    if len(stderr) > 0 { return nil, errors.New(string(stderr)) }
    // fmt.Printf("%s\n", stdout)

    location, _ := time.LoadLocation("Local")
    timestring := strings.TrimSpace(string(stdout))
    pinfo.StartTime, error = time.ParseInLocation(time.ANSIC, timestring, location)
    if error != nil { return nil, error }

    //  Do command line & args --

    stdout, stderr, error = RunShellCommand("ps", []string{"-ocomm=", "-p", pidStr}, nil)
    if error != nil { return nil, error }
    if len(stderr) > 0 { return nil, errors.New(string(stderr)) }
    // fmt.Printf("%s\n", stdout)
    pinfo.Command = strings.TrimSpace(string(stdout))

    stdout, stderr, error = RunShellCommand("ps", []string{"-ocommand=", "-p", pidStr}, nil)
    if error != nil { return nil, error }
    if len(stderr) > 0 { return nil, errors.New(string(stderr)) }
    // fmt.Printf("%s\n", stdout)
    pinfo.CommandLine = strings.TrimSpace(string(stdout))

    return pinfo, nil
}

