//  claim-invite  -  Claim an invite link.
//
//  E.B.Smith  -  November, 2014


package main


import (
	"database/sql"	
 	_ "github.com/go-sql-driver/mysql"
	)


var globalDatabase *sql.DB = nil


func connectDatabase() error {
	//	Make a connection --
	var error error = nil
	globalDatabase, error = sql.Open("mysql", "RelcyUserAdmin:RelcyUserAdmin@/RelcyUserDatabase")

	if error == nil {
		ZLog(ZLogDebug, "Database opened.")
	} else {
		globalDatabase = nil
		ZLog(ZLogError, "Can't open database connection: %v.", error);
		return error
		}

	return error
	}


func disconnectDatabase() {
	if  globalDatabase != nil {
		globalDatabase.Close()
		globalDatabase = nil
		}
	}
