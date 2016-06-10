//  RunShellCommand_test  -  RunShellCommand tests
//
//  E.B.Smith  -  March, 2015


package Util


import (
    "os"
    "fmt"
    "path"
    "os/user"
    "testing"
)



//----------------------------------------------------------------------------------------
//                                                                            TestHomePath
//----------------------------------------------------------------------------------------


func TestHomePath(t *testing.T) {
    u, _ := user.Current()
    h := u.HomeDir
    r := HomePath()
    if h != r {
        t.Errorf("Got %s but want %s.", r, h)
    }
}



//----------------------------------------------------------------------------------------
//                                                                         TestCleanupPath
//----------------------------------------------------------------------------------------


func TestCleanupPath(t *testing.T) {
    homedir := HomePath()
    currdir, _ := os.Getwd()
    tests := []struct {
        testin, testout string
    }{
        { "bob",                currdir+"/bob" },
        { "~/bob",              homedir+"/bob" },
        { "~bob",               homedir+"/bob" },
        { ".",                  currdir },
        { "./",                 currdir },
        { "../cc/blab",         path.Dir(currdir)+"/cc/blab" },
        { "../../Documents",    path.Dir(path.Dir(currdir))+"/Documents" },
        { "Test/../What/",      currdir+"/What" },
        { "/etc/var",           "/etc/var" },
        { "/etc/var/",          "/etc/var" },
    }

    for _, test := range tests {
        result := AbsolutePath(test.testin)
        if false { fmt.Printf("%s\t\t%s\n", result, test.testout) }
        if result != test.testout {
            t.Errorf("Got %s but want %s.", result, test.testout)
        }
    }
}


