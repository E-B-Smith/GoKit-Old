//  psql  -  A go Postgres interface.
//
//  E.B.Smith  -  November, 2014


package psql


import (
    "io"
    "fmt"
    "sync"
    "bufio"
    "errors"
    "strings"
    "strconv"
    "os/exec"
    "net/url"
    "database/sql"
    _ "github.com/lib/pq"
    "violent.blue/go/log"
    )


var PGCTLPath string    = ""
var PSQLPath string     = ""
var PSQLDataPath string = ""
var DB *sql.DB          = nil

var Databasename = ""
var Host         = "localhost"
var Username     = "postgres"
var Password     = ""
var Port         = 5432


func ConnectDatabase(databaseURI string) error {
    //
    //  Start the database --
    //

    if databaseURI != "" {
        u, error := url.Parse(databaseURI)
        if error != nil {
            return error
        } else {
        if u == nil {
            return errors.New("Invalid database URI")
            }}
        log.Debug("%s:\n%v", databaseURI, u)

        if u.Scheme == "db" || u.Scheme == "psql" || u.Scheme == "sql" {
        } else {
            log.Error("Invalid database scheme '%s'", u.Scheme)
            return errors.New("Invalid scheme")
            }

        i := strings.IndexRune(u.Host, ':')
        if i >= 0 {
            Host = u.Host[0:i]
            Port, _ = strconv.Atoi(u.Host[i+1:])
            }
        if Port <= 0 { Port = 5432 }
        if u.User == nil {
            Username = ""
            Password = ""
        } else {
            Username = u.User.Username()
            Password, _ = u.User.Password()
            }
        Databasename = u.Path
        if len(Databasename) > 1 && Databasename[0:1] == "/" { Databasename = Databasename[1:] }
        log.Debug("Host: %s Port: %d User: %s Pass: %s Databasename: %s.", Host, Port, Username, Password, Databasename)
        }

    //  Find postgres --

    var error error
    PGCTLPath, error = exec.LookPath("pg_ctl")
    if error != nil {
        log.Error("Can't find Postgres 'pg_ctl': %v.", error)
        return error
        }
    log.Debug("   Found pg_ctl: %v.", PGCTLPath)
    var command *exec.Cmd
    if len(PSQLDataPath) > 0 {
        log.Debug("Using data path: %v.", PSQLDataPath)
        command = exec.Command(PGCTLPath, "status", "-D",  PSQLDataPath)
    } else {
        log.Debug("Using default datapath.")
        command = exec.Command(PGCTLPath, "status")
    }
    error = command.Run()
    if command.ProcessState.Sys() == 3 {
        log.Debug("Starting Postgres")
        if len(PSQLDataPath) > 0 {
           command = exec.Command(PGCTLPath, "start", "-w", "-s", "-D", PSQLDataPath)
        } else {
           command = exec.Command(PGCTLPath, "start", "-w", "-s")
        }
        error = command.Run()
        if error != nil {
            log.Error("Can't start Postgress: %v.", error)
            return error
        }
    } else {
        log.Debug("Postgres is already started.")
    }


    //
    //  Find psql command line utility and connect --
    //


    //  Find psql --

    PSQLPath, error = exec.LookPath("psql")
    if error != nil {
        log.Error("Can't find Postgres 'psql': %v.", error)
        return error
        }
    log.Debug("psqlpath: %v.", PSQLPath)

    //  Make a connection --

    connectString :=
        fmt.Sprintf("host=%s port=%d  dbname=%s user=%s password=%s sslmode=disable",
                     Host, Port, Databasename, Username, Password)
    log.Debug("Connection string: %s.", connectString)
    DB, error = sql.Open("postgres", connectString)
    if error != nil {
        DB = nil
        log.Error("Error: Can't open database connection: %v.", error);
        return error
    }

    //  Get our settings --
    //  select setting from pg_settings where name = 'port';
    rows, error := DB.Query("select current_user, inet_server_addr(), inet_server_port(), current_database(), current_schema;")
    if error != nil {
        log.Error("Error querying database config: %v.", error)
        return error
    } else {
        defer rows.Close()
        var (user string; host string; port int; database string; schema string)
        for rows.Next() {
            rows.Scan(&user, &host, &port, &database, &schema)
            log.Debug("Connected to database psql:%s@%s:%d/%s (%s).", user, host, port, database, schema)
        }
    }

    return nil
}


func DisconnectDatabase() {
    if  DB != nil {
        DB.Close()
        DB = nil
        Host = "localhost"
        Port = 5432
        Databasename = "postgres"
        Username = "postgres"
    }
}


func RunScript(script string) error {

    //
    //  Run an SQL script that is stored as a resource --
    //

    var error error
    psqlOptions := [] string {
        "-h", "localhost",
        "-X", "-q",
        "-v", "ON_ERROR_STOP=1",
        "--pset", "pager=off",
    }
    command := exec.Command(PSQLPath, psqlOptions...)
    command.Env = append(command.Env, "PGOPTIONS=-c client_min_messages=WARNING")
    commandpipe, error := command.StdinPipe()
    if error != nil {
        log.Error("Can't open pipe: %v", error)
        return error
    }

    var errorpipe *io.PipeReader;
    errorpipe, command.Stderr = io.Pipe()

    error = command.Start()
    if error != nil {
        log.Error("Error running psql: %v.", error)
        return error
    }

    commandpipe.Write([]byte(script))
    commandpipe.Close()

    var waiter sync.WaitGroup
    waiter.Add(1)
    go func() {
        scanner := bufio.NewScanner(errorpipe)
        for scanner.Scan() {
            log.Error("%v.", scanner.Text())
        }
        waiter.Done()
    } ()

    error = command.Wait()
    errorpipe.Close()
    waiter.Wait()

    if error != nil {
        log.Error("Script %v.", error)
        return error
    }

    return nil
}


func StringFromArray(ary []string) string {
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

func ArrayFromString(s *string) []string {
    if s == nil { return *new([]string) }

    str := strings.Trim(*s, "{}")
    a := make([]string, 0, 10)
    for _, ss := range strings.Split(str, ",") {
        a = append(a, ss)
    }
    return a
}


func RunScript2(script string) (standardOut []byte, standardError []byte, error error) {
    //
    //  Execute an SQL script --
    //

    psqlOptions := [] string {
        "-X", "-q",
        "-v", "ON_ERROR_STOP=1",
        "--pset", "pager=off",
    }
    if Host == "" {
        psqlOptions = append(psqlOptions, "-h", "localhost")
    } else {
        psqlOptions = append(psqlOptions, "-h", Host)
    }
    psqlOptions = append(psqlOptions, fmt.Sprintf("--port=%d", Port))
    psqlOptions = append(psqlOptions, Databasename, Username)

    command := exec.Command(PSQLPath, psqlOptions...)
    command.Env = append(command.Env, "PGOPTIONS=-c client_min_messages=WARNING")
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


