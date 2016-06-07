//  Paths  -  Path functions.
//
//  E.B.Smith  -  March, 2015


package Util


import (
    "os"
    "os/user"
    "path"
    "strings"
    "path/filepath"
)


//  Returns the user's home directory
func HomePath() string {
    homepath := ""
    u, error := user.Current()
    if error == nil {
        homepath = u.HomeDir
    } else {
        homepath = os.Getenv("HOME")
    }
    return homepath
}


//  Like the regular AbsolutePath but adds home directory if indicated.
func AbsolutePath(filename string) string {
    filename = strings.TrimSpace(filename)
    if  filepath.HasPrefix(filename, "~") {
        filename = strings.TrimPrefix(filename, "~")
        filename = path.Join(HomePath(), filename)
    }
    if ! path.IsAbs(filename) {
        s, _ := os.Getwd()
        filename = path.Join(s, filename)
    }
    filename = path.Clean(filename)
    return filename
}

