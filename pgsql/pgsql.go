//  pgsql  -  A go Postgres interface.
//
//  E.B.Smith  -  November, 2014


package pgsql


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
    "violent.blue/GoKit/Log"
)



//----------------------------------------------------------------------------------------
//                                                                                   pgsql
//----------------------------------------------------------------------------------------


type SSLType int;
const (
    SSLTypeDisable = iota   //  ?sslmode=disable
    SSLTypeRequire          //  require
    SSLTypeVerifyCA         //  verify-ca
    SSLTypeVerifyFull       //  verify-full
)


type PGSQL struct {
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
    UseSSL          bool
}


func DefaultValue() PGSQL {
    pgsql := PGSQL {
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
        UseSSL:         false,
    }
    return pgsql
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
//                                                                                 Helpers
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


type RowScanner interface {
    Scan(dest ...interface{}) error
}


func RowsUpdated(result sql.Result) int64 {
    var rowsUpdated int64 = 0
    if result != nil { rowsUpdated, _ = result.RowsAffected() }
    return rowsUpdated
}


func UpdateResultError(result sql.Result, error error) error {
    if error != nil {
        return error
    }
    if result == nil {
        return fmt.Errorf("No rows updated")
    }
    var rowsUpdated int64 = 0
    rowsUpdated, error = result.RowsAffected()
    if error != nil {
        return error
    }
    if rowsUpdated == 0 {
        return fmt.Errorf("No rows updated")
    }
    return nil
}


//----------------------------------------------------------------------------------------
//                                                                         ConnectDatabase
//----------------------------------------------------------------------------------------


func ConnectDatabase(databaseURI string) (psql *PGSQL, error error) {
    //
    //  Connect to the database --
    //

    //  Parse a URI like:
    //  psql://happylabsadmin:happylabsadmin@localhost:5432/happylabsdatabase

    psqlValue := DefaultValue()
    psql = &psqlValue

    if databaseURI != "" {
        u, error := url.Parse(databaseURI)
        if error != nil {
            return nil, error
        } else if u == nil {
            return nil, errors.New("Invalid database URI")
        }
        Log.Debugf("URI: %s. Parsed: %+v.", databaseURI, u)

        if u.Scheme == "db"   ||
           u.Scheme == "psql" ||
           u.Scheme == "sql"  ||
           u.Scheme == "postgres" {
        } else {
            Log.Errorf("Invalid database scheme '%s'", u.Scheme)
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
        Log.Debugf("Host: %s Port: %d User: %s Pass: %s Databasename: %s.",
            psql.Host, psql.Port, psql.Username, psql.password, psql.Databasename)
    }

    //  Find postgres --

    psql.PGCTLPath, error = exec.LookPath("pg_ctl")
    if error != nil {
        Log.Errorf("Can't find Postgres 'pg_ctl': %v.", error)
        return nil, error
    }
    Log.Debugf("   Found pg_ctl: %v.", psql.PGCTLPath)

    //  Is postgres running?

    var command *exec.Cmd
    if len(psql.PSQLDataPath) > 0 {
        Log.Debugf("Using data path: %v.", psql.PSQLDataPath)
        command = exec.Command(psql.PGCTLPath, "status", "-D",  psql.PSQLDataPath)
    } else {
        Log.Debugf("Using default datapath.")
        command = exec.Command(psql.PGCTLPath, "status")
    }
    error = command.Run()
    if ExitCodeFromProcessState(command.ProcessState) == 3 {
        Log.Debugf("Starting Postgres")
        if len(psql.PSQLDataPath) > 0 {
           command = exec.Command(psql.PGCTLPath, "start", "-w", "-D", psql.PSQLDataPath)
        } else {
           command = exec.Command(psql.PGCTLPath, "start", "-w")
        }
        error = command.Run()
        if error != nil {
            Log.Errorf("Can't start Postgress: %v.", error)
            return nil, error
        }
    } else {
        Log.Debugf("Postgres is already started.")
    }

    //
    //  Find psql command line utility and connect --
    //

    //  Find psql --

    psql.PSQLPath, error = exec.LookPath("psql")
    if error != nil {
        Log.Errorf("Can't find Postgres 'psql': %v.", error)
        return nil, error
    }
    Log.Debugf("psqlpath: %v.", psql.PSQLPath)

    //  Make a connection --

    connectString := fmt.Sprintf("host=%s port=%d", psql.Host, psql.Port)
    if psql.Databasename != "" {
       connectString += " dbname="+psql.Databasename
    }
    if psql.Username != "" {
        connectString += " user="+psql.Username
    }
    if psql.password != "" {
        connectString += " password="+psql.password
    }
    if psql.UseSSL {
        connectString += " sslmode=enable"
    } else {
        connectString += " sslmode=disable"
    }

    Log.Debugf("Connection string: %s.", connectString)
    //connectString = databaseURI //  eDebug
    psql.DB, error = sql.Open("postgres", connectString)
    if error != nil {
        psql.DB = nil
        Log.Errorf("Error: Can't open database connection: %v.", error);
        return nil, error
    }

    //  Get our settings --
    //  select setting from pg_settings where name = 'port';

    rows, error := psql.DB.Query("select current_user, inet_server_addr(), inet_server_port(), current_database(), current_schema;")
    defer CloseRows(rows)
    if error != nil {
        Log.Errorf("Error querying database config: %v.", error)
        return nil, error
    } else {
        var (user string; host string; port int; database string; schema string)
        for rows.Next() {
            rows.Scan(&user, &host, &port, &database, &schema)
            Log.Debugf("Connected to database psql:%s@%s:%d/%s (%s).", user, host, port, database, schema)
        }
    }

    return psql, nil
}


func (psql *PGSQL) DisconnectDatabase() {
    if  psql.DB != nil {
        psql.DB.Close()
        *psql = DefaultValue()
    }
}


//----------------------------------------------------------------------------------------
//                                                                            RunSQLScript
//----------------------------------------------------------------------------------------


func (psql *PGSQL) RunSQLScript(script string) (standardOut []byte, standardError []byte, error error) {
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
        Log.Errorf("Can't open StdIn pipe: %v.", error)
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
        for error == nil {
            standardOut = append(standardOut, buffer[:count]...)
            count, error = reader.Read(buffer)
            //Log.Debugf("Read %d bytes.  Error: %v.", count, error)
        }
        waiter.Done()
    } ()

    _, error = stdinpipe.Write([]byte(script))
    if error != nil {
        Log.Errorf("Error writing stdin: %v.", error)
        return standardOut, standardError, error
    }

    error = command.Start()
    if error != nil {
        Log.Errorf("Error starting psql: %v.", error)
        return standardOut, standardError, error
    }

    stdinpipe.Close()
    error = command.Wait()
    errorpipe.Close()
    pipeout.Close()
    waiter.Wait()

    if error != nil {
        Log.Errorf("Script wait error: %v.", error)
        return standardOut, standardError, error
    }

    return standardOut, standardError, nil
}



