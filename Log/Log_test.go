//  Log  -  Simple logging.
//
//  E.B.Smith  -  November, 2014


package Log


import (
//    "os"
    "time"
    "testing"
    "path/filepath"
)


func TestLogRotation(t *testing.T) {

    LogTeeStderr = true
    LogLevel = LevelAll
    LogRotationInterval = time.Second * 5
    SetFilename("Log/TestLog.log")

    Infof("Starting test.")
    for i := 0; i < 10*5; i++ {
        Infof("Message %d.", i)
        time.Sleep(time.Second)
    }

    logfiles, _ := filepath.Glob("Log/TestLog*")
    if len(logfiles) != 8 {
        t.Errorf("Expected 8 files, found %d.", len(logfiles))
    }

    // for _, file := range logfiles {
    //     error = os.Remove(file)
    //     if error != nil {
    //         t.Errorf("Can't remove log file '%s': %v.", file, error)
    //     }
    // }
}
