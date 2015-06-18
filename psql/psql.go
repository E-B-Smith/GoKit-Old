//  psql  -  A go Postgres interface.
//
//  E.B.Smith  -  November, 2014


package psql


import (
    "io"
    "os"
    "fmt"
    "sync"
    "time"
    "bufio"
    "errors"
    "syscall"
    "strings"
    "strconv"
    "os/exec"
    "net/url"
    "database/sql"
    "github.com/lib/pq"
    "violent.blue/go/log"
    )



//----------------------------------------------------------------------------------------
//                                                                                    psql
//----------------------------------------------------------------------------------------


type PSQL struct {
    PGCTLPath       string
    PSQLPath        string
    PSQLDataPath    string
    DB              *sql.DB
    Databasename    string
    Host            string
    Username        string
    password        string
    Port            int
    infiniteTimeEnabled  bool
}


func DefaultValue() PSQL {
    psql := PSQL {
        PGCTLPath:      "",
        PSQLPath:       "",
        PSQLDataPath:   "",
        DB:             nil,
        Databasename:   "",
        Host:           "localhost",
        Username:       "postgres",
        password:       "",
        Port:           5432,
        infiniteTimeEnabled: false,
    }
    return psql
}



//----------------------------------------------------------------------------------------
//                                                                      EnableInfiniteTime
//----------------------------------------------------------------------------------------


var NegativeInfinityTime time.Time = time.Date(1500, time.January, 1, 0, 0, 0, 0, time.UTC)
var PositiveInfinityTime time.Time = time.Date(2500, time.January, 1, 0, 0, 0, 0, time.UTC)
var infiniteTimeEnabled bool = false

func EnableInfiniteTime() {
    if !infiniteTimeEnabled {               //  eDebug: Threading
        infiniteTimeEnabled = true
        pq.EnableInfinityTs(NegativeInfinityTime, PositiveInfinityTime)
    }
}



//----------------------------------------------------------------------------------------
//                                                                         ConnectDatabase
//----------------------------------------------------------------------------------------


func ExitCodeFromProcessState(ps *os.ProcessState) int {
    if ps == nil { return 1 }
    if status, ok := ps.Sys().(syscall.WaitStatus); ok {
        return status.ExitStatus()
    }
    return 1
}


func CloseRows(rows *sql.Rows)  {
    if rows != nil {
        rows.Close();
    }
}


func ConnectDatabase(databaseURI string) (psql *PSQL, error error) {
    //
    //  Start the database --
    //

    psql = new(PSQL)

    if databaseURI != "" {
        u, error := url.Parse(databaseURI)
        if error != nil {
            return nil, error
        } else if u == nil {
            return nil, errors.New("Invalid database URI")
        }
        log.Debug("%s:\n%v", databaseURI, u)

        if u.Scheme == "db" || u.Scheme == "psql" || u.Scheme == "sql" {
        } else {
            log.Error("Invalid database scheme '%s'", u.Scheme)
            return nil, errors.New("Invalid scheme")
        }

        i := strings.IndexRune(u.Host, ':')
        if i >= 0 {
            psql.Host = u.Host[0:i]
            psql.Port, _ = strconv.Atoi(u.Host[i+1:])
        }
        if psql.Port <= 0 { psql.Port = 5432 }
        if u.User == nil {
            psql.Username = ""
            psql.password = ""
        } else {
            psql.Username = u.User.Username()
            psql.password, _ = u.User.Password()
        }
        psql.Databasename = u.Path
        if len(psql.Databasename) > 1 && psql.Databasename[0:1] == "/" {
            psql.Databasename = psql.Databasename[1:]
        }
        log.Debug("Host: %s Port: %d User: %s Pass: %s Databasename: %s.",
            psql.Host, psql.Port, psql.Username, psql.password, psql.Databasename)
    }

    //  Find postgres --

    psql.PGCTLPath, error = exec.LookPath("pg_ctl")
    if error != nil {
        log.Error("Can't find Postgres 'pg_ctl': %v.", error)
        return nil, error
    }
    log.Debug("   Found pg_ctl: %v.", psql.PGCTLPath)

    //  Is postgres running?

    var command *exec.Cmd
    if len(psql.PSQLDataPath) > 0 {
        log.Debug("Using data path: %v.", psql.PSQLDataPath)
        command = exec.Command(psql.PGCTLPath, "status", "-D",  psql.PSQLDataPath)
    } else {
        log.Debug("Using default datapath.")
        command = exec.Command(psql.PGCTLPath, "status")
    }
    error = command.Run()
    if ExitCodeFromProcessState(command.ProcessState) == 3 {
        log.Debug("Starting Postgres")
        if len(psql.PSQLDataPath) > 0 {
           command = exec.Command(psql.PGCTLPath, "start", "-w", "-D", psql.PSQLDataPath)
        } else {
           command = exec.Command(psql.PGCTLPath, "start", "-w")
        }
        error = command.Run()
        if error != nil {
            log.Error("Can't start Postgress: %v.", error)
            return nil, error
        }
    } else {
        log.Debug("Postgres is already started.")
    }

    //
    //  Find psql command line utility and connect --
    //

    //  Find psql --

    psql.PSQLPath, error = exec.LookPath("psql")
    if error != nil {
        log.Error("Can't find Postgres 'psql': %v.", error)
        return nil, error
    }
    log.Debug("psqlpath: %v.", psql.PSQLPath)

    //  Make a connection --

    connectString :=
        fmt.Sprintf("host=%s port=%d  dbname=%s user=%s password=%s sslmode=disable",
                     psql.Host, psql.Port, psql.Databasename, psql.Username, psql.password)
    log.Debug("Connection string: %s.", connectString)
    psql.DB, error = sql.Open("postgres", connectString)
    if error != nil {
        psql.DB = nil
        log.Error("Error: Can't open database connection: %v.", error);
        return nil, error
    }

    //  Get our settings --
    //  select setting from pg_settings where name = 'port';

    rows, error := psql.DB.Query("select current_user, inet_server_addr(), inet_server_port(), current_database(), current_schema;")
    defer CloseRows(rows)
    if error != nil {
        log.Error("Error querying database config: %v.", error)
        return nil, error
    } else {
        var (user string; host string; port int; database string; schema string)
        for rows.Next() {
            rows.Scan(&user, &host, &port, &database, &schema)
            log.Debug("Connected to database psql:%s@%s:%d/%s (%s).", user, host, port, database, schema)
        }
    }

    return psql, nil
}


func (psql *PSQL) DisconnectDatabase() {
    if  psql.DB != nil {
        psql.DB.Close()
        *psql = DefaultValue()
    }
}



//----------------------------------------------------------------------------------------
//                                                                                  Arrays
//----------------------------------------------------------------------------------------


func StringFromStringArray(ary []string) string {
    if len(ary) == 0 {
        return "{}"
    }

    var result string = "{"+ary[0];
    for i:=1; i < len(ary); i++ {
        result += ","+ary[i]
    }
    result += "}"
    return result
}


func StringArrayFromString(s *string) []string {
    if s == nil { return *new([]string) }

    str := strings.Trim(*s, "{}")
    a := make([]string, 0, 10)
    for _, ss := range strings.Split(str, ",") {
        a = append(a, ss)
    }
    return a
}


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
        if error != nil { a = append(a, int32(i)) }
    }
    return a
}



//----------------------------------------------------------------------------------------
//                                                                            RunSQLScript
//----------------------------------------------------------------------------------------


func (psql *PSQL) RunSQLScript(script string) (standardOut []byte, standardError []byte, error error) {
    //
    //  Execute an SQL script --
    //

    psqlOptions := [] string {
        "-X", "-q",
        "-v", "ON_ERROR_STOP=1",
        "--pset", "pager=off",
    }
    if psql.Host == "" {
        psqlOptions = append(psqlOptions, "-h", "localhost")
    } else {
        psqlOptions = append(psqlOptions, "-h", psql.Host)
    }
    psqlOptions = append(psqlOptions, fmt.Sprintf("--port=%d", psql.Port))
    psqlOptions = append(psqlOptions, psql.Databasename, psql.Username)

    command := exec.Command(psql.PSQLPath, psqlOptions...)
    command.Env = append(command.Env, "PGOPTIONS=-c client_min_messages=WARNING")
    if len(psql.password) > 0 {
        command.Env = append(command.Env, fmt.Sprintf("PGPASSWORD=%s", psql.password))
    }
    stdinpipe, error := command.StdinPipe()
    if error != nil {
        log.Error("Can't open StdIn pipe: %v.", error)
        return standardOut, standardError, error
    }

    var waiter sync.WaitGroup
    waiter.Add(2)

    var errorpipe *io.PipeReader
    errorpipe, command.Stderr = io.Pipe()
    go func() {
        buffer := make([]byte, 512)
        reader := bufio.NewReader(errorpipe)
        count, error := reader.Read(buffer)
        //log.Debug("Read %d bytes.  Error: %v.", count, error)
        for error == nil {
            standardError = append(standardError, buffer[:count]...)
            count, error = reader.Read(buffer)
            //log.Debug("Read %d bytes.  Error: %v.", count, error)
        }
        waiter.Done()
    } ()

    var pipeout *io.PipeReader
    pipeout, command.Stdout = io.Pipe()
    go func () {
        buffer := make([]byte, 512)
        reader := bufio.NewReader(pipeout)
        count, error := reader.Read(buffer)
        //log.Debug("Read %d bytes.  Error: %v.", count, error)
        for error == nil {
            standardOut = append(standardOut, buffer[:count]...)
            count, error = reader.Read(buffer)
            //log.Debug("Read %d bytes.  Error: %v.", count, error)
        }
        waiter.Done()
    } ()

    _, error = stdinpipe.Write([]byte(script))
    if error != nil {
        log.Error("Error writing stdin: %v.", error)
        return standardOut, standardError, error
    }

    error = command.Start()
    if error != nil {
        log.Error("Error starting psql: %v.", error)
        return standardOut, standardError, error
    }

    stdinpipe.Close()
    error = command.Wait()
    errorpipe.Close()
    pipeout.Close()
    waiter.Wait()

    if error != nil {
        log.Error("Script wait error: %v.", error)
        return standardOut, standardError, error
    }

    return standardOut, standardError, nil
}



