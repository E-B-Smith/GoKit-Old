//  Scheduler_test.go  -  Test the scheduler.
//
//  E.B.Smith  -  March, 2015


package Util


import (
    "sync"
    "time"
    "testing"
    "violent.blue/GoKit/Log"
)


//----------------------------------------------------------------------------------------
//                                                                           TestScheduler
//----------------------------------------------------------------------------------------


func absDur(i time.Duration) time.Duration {
    if i < 0 { return -1 * i }
    return i
}

func TestScheduler(test *testing.T) {

    Log.LogLevel = Log.LogLevelAll

    lock := new(sync.Mutex)
    timemap := make(map[int]time.Time)
    var tests, errors int

    testFun := func(idx int, i int) {

        lock.Lock()
        defer lock.Unlock()

        if exp, ok := timemap[idx]; ok {
            now := time.Now()
            diff := absDur(now.Sub(exp))
            if  diff > 6*time.Millisecond {
                test.Errorf("To long! For %d %d:\nNow: %v\nExp: %v\nDiff: %v.\n", idx, i, now, exp, diff)
                errors++
            }
        }
        timemap[idx] = time.Now().Add(time.Duration(i) * time.Millisecond)
        tests++
    }


    StartScheduler()

    ScheduleTask(10 * time.Millisecond, func() { testFun( 1, 10) })
    ScheduleTask(20 * time.Millisecond, func() { testFun( 2, 20) })
    ScheduleTask(30 * time.Millisecond, func() { testFun( 3, 30) })
    ScheduleTask(40 * time.Millisecond, func() { testFun( 4, 40) })
    ScheduleTask(50 * time.Millisecond, func() { testFun( 5, 50) })
    ScheduleTask(60 * time.Millisecond, func() { testFun( 6, 60) })
    ScheduleTask(70 * time.Millisecond, func() { testFun( 7, 70) })
    ScheduleTask(80 * time.Millisecond, func() { testFun( 8, 80) })
    ScheduleTask(90 * time.Millisecond, func() { testFun( 9, 90) })

    ScheduleTask(90 * time.Millisecond, func() { testFun(10, 90) })
    ScheduleTask(80 * time.Millisecond, func() { testFun(11, 80) })
    ScheduleTask(70 * time.Millisecond, func() { testFun(12, 70) })
    ScheduleTask(60 * time.Millisecond, func() { testFun(13, 60) })
    ScheduleTask(50 * time.Millisecond, func() { testFun(14, 50) })
    ScheduleTask(40 * time.Millisecond, func() { testFun(15, 40) })
    ScheduleTask(30 * time.Millisecond, func() { testFun(16, 30) })
    ScheduleTask(20 * time.Millisecond, func() { testFun(17, 20) })
    ScheduleTask(10 * time.Millisecond, func() { testFun(18, 10) })

    ScheduleTask(40 * time.Millisecond, func() { testFun(19, 40) })
    ScheduleTask(40 * time.Millisecond, func() { testFun(20, 40) })
    ScheduleTask(40 * time.Millisecond, func() { testFun(21, 40) })
    ScheduleTask(40 * time.Millisecond, func() { testFun(22, 40) })
    ScheduleTask(40 * time.Millisecond, func() { testFun(23, 40) })

    ScheduleTask(20 * time.Millisecond, func() { testFun(24, 20) })
    ScheduleTask(20 * time.Millisecond, func() { testFun(25, 20) })

    time.Sleep(10 * time.Second)

    StopScheduler()
    Log.Infof("Tests: %d Errors: %d.", tests, errors)
}

