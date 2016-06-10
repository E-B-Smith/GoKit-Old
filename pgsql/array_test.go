//  pgsql.array_test  -  A go Postgres interface for handling postgres arrays.
//
//  E.B.Smith  -  November, 2014


package pgsql


import (
    "testing"
    "database/sql"
    "violent.blue/GoKit/Log"
)


func StringPtr(s string) *string {
    return &s
}


func StringArraysEqual(a []string, b []string) bool {
    if len(a) != len(b) { return false }

    for i := 0; i < len(a); i++ {
        if a[i] != b[i] { return false }
    }

    return true
}


func Test_ReadStringArrayFromDB(t *testing.T) {
    Log.LogLevel = Log.LogLevelAll

    sa := StringArrayFromNullString(sql.NullString {})
    if len(sa) != 0 {
        t.Error("Expected string with len 0 from null string.")
    }

    sa = StringArrayFromNullString(sql.NullString { String: `{"1","2","3"}`, Valid: false })
    if len(sa) != 0 {
        t.Error("Expected string with len 0 from null string.")
    }

    sa = StringArrayFromNullString(sql.NullString { String: `{}`, Valid: true })
    if len(sa) != 0 {
        t.Errorf("Expected array with len 0 from empty input array. Got: (length %d) %v.", len(sa), sa)
    }

    var ta []string

    ta = []string { "1" }
    sa = StringArrayFromNullString(sql.NullString { String: `{"1"}`, Valid: true })
    if ! StringArraysEqual(sa, ta) {
        t.Errorf("Expected: %v\nGot: %v", ta, sa)
    }

    ta = []string { "1", "2", "3" }
    sa = StringArrayFromNullString(sql.NullString { String: `{"1","2","3"}`, Valid: true })
    if ! StringArraysEqual(sa, ta) {
        t.Errorf("Expected: %v Got: %v", ta, sa)
    }

    ta = []string { "Consulting offer", "New venture", "Sale leads" }
    sa = StringArrayFromNullString(sql.NullString { String: `{"Consulting offer","New venture","Sale leads"}`, Valid: true })
    if ! StringArraysEqual(sa, ta) {
        t.Errorf("Expected: %v Got: %v", ta, sa)
    }

    //  {"Consulting offer","New venture","Sale leads"}
}

