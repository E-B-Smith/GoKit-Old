//  psql  -  A go Postgres interface.
//
//  E.B.Smith  -  November, 2014


package psql


import (
	"io"
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


var PGCTLPath string 	= ""
var PSQLPath string 	= ""
var PSQLDataPath string = ""
var DB *sql.DB			= nil

var Name		= ""
var Host        = "localhost"
var User		= "postgres"
var Password	= ""
var Port		= 5432


func ConnectDatabase(databaseURI string) error {
	//
	//	Start the database --
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
            User = ""
            Password = ""
        } else {
            User = u.User.Username()
            Password, _ = u.User.Password()
            }
        Name = u.Path
        if len(Name) > 1 && Name[0:1] == "/" { Name = Name[1:] }
        log.Debug("Host: %s Port: %d User: %s Pass: %s Name: %s.", Host, Port, User, Password, Name)
        }

	//	Find postgres --

	var error error
	PGCTLPath, error = exec.LookPath("pg_ctl")
	if error != nil {
		log.Error("Can't find Postgres 'pg_ctl': %v.", error)
		return error
		}
	log.Debug(" Found pg_ctl: %v.", PGCTLPath)
	log.Debug("Database Data: %v.", PSQLDataPath)

	command := exec.Command(PGCTLPath, "status", "-D", 	PSQLDataPath)
	error = command.Run()
	if command.ProcessState.Sys() == 3 {
		log.Debug("Starting Postgres")
		command = exec.Command(PGCTLPath, "start", "-w", "-s", "-D", PSQLDataPath)
		error = command.Run()
		if error != nil {
			log.Error("Can't start Postgress: %v.", error)
			return error
			}
	} else {
		log.Debug("Postgres is already started.")
		}


	//
	//	Find psql command line utility and connect --
	//


	//	Find psql --

	PSQLPath, error = exec.LookPath("psql")
	if error != nil {
		log.Error("Can't find Postgres 'psql': %v.", error)
		return error
		}
	log.Debug("psqlpath: %v.", PSQLPath)

	//	Make a connection --



	DB, error = sql.Open("postgres", "user=Edward dbname=Edward sslmode=disable")
	if error != nil {
		DB = nil
		log.Error("Error: Can't open database connection: %v.", error);
		return error
		}

	return nil
	}


func DisconnectDatabase() {
	if  DB != nil {
		DB.Close()
		DB = nil
		Host = "localhost"
		Port = 5432
		Name = "postgres"
		User = "postgres"
		}
	}


func RunScript(script string) error {

	//
	//	Run an SQL script that is stored as a resource --
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


