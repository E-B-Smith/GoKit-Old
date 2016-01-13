//  RunShellCommand_test  -  RunShellCommand tests
//
//  E.B.Smith  -  March, 2015


package Util


import (
    "os"
    "fmt"
    "testing"
)



//----------------------------------------------------------------------------------------
//                                                                     TestRunShellCommand
//----------------------------------------------------------------------------------------


func TestRunShellCommand(t *testing.T) {
    stdin := "echo 'It works.'; echo 'StdErr' >&2;"

    stdout, stderr, error := RunShellCommand("bash", nil, []byte(stdin))

    if error != nil || string(stdout) != "It works.\n" || string(stderr) != "StdErr\n" {
        t.Errorf("Got error: %v stdout: %s stderr: %s.\n", error, string(stdout), string(stderr))
    }
    //fmt.Printf("%v:\t'%s'\t'%s'\n", error, stdout, stderr)

    stdout, stderr, error = RunShellCommand("bbbbash", nil, []byte(stdin))
    if error == nil || stdout != nil || stderr != nil {
        t.Errorf("Got error: %v stdout: %s stderr: %s.\n", error, string(stdout), string(stderr))
    }
    //fmt.Printf("%v:\t'%s'\t'%s'\n", error, stdout, stderr)

    stdout, stderr, error = RunShellCommand("/bin/echo", []string{"'Here'"}, nil)
    if error != nil || string(stdout) != "'Here'\n" || string(stderr) != "" {
        t.Errorf("Got error: %v stdout: %s stderr: %s.\n", error, string(stdout), string(stderr))
    }

}



//----------------------------------------------------------------------------------------
//                                                                          GetProcessInfo
//----------------------------------------------------------------------------------------


func TestGetProcessInfo(t *testing.T) {
    pinfo, error := GetProcessInfo(os.Getpid())
    if error != nil { t.Errorf("Got error %v.", error) }
    fmt.Printf("%+v\n", pinfo)
}

