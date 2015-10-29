//  Paths  -  Path functions.
//
//  E.B.Smith  -  March, 2015


package ServerUtil


import (
    "os"
    "os/user"
    "path"
    "strings"
    "path/filepath"
)


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


func CleanupPath(filename string) string {
    if  filepath.HasPrefix(filename, "~") {
        home := HomePath() + "/"
        filename = strings.Replace(filename, "~", home, 1)
    }
    if ! path.IsAbs(filename) {
        s, _ := os.Getwd()
        filename = path.Join(s, filename)
    }
    filename = path.Clean(filename)
    return filename
}

