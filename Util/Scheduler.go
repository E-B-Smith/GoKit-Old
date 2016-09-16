

//----------------------------------------------------------------------------------------
//
//                                                               GoKit/Util : Scheduler.go
//                                Periodicaly runs scheduled tasks.  Naive task scheduler.
//
//                                                                 E.B. Smith, August 2016
//                        -©- Copyright © 2014-2016 Edward Smith, all rights reserved. -©-
//
//----------------------------------------------------------------------------------------


package Util


import (
    "time"
    "container/heap"
    "violent.blue/GoKit/Log"
)


//----------------------------------------------------------------------------------------
//
//                                                                               Scheduler
//
//----------------------------------------------------------------------------------------


//----------------------------------------------------------------------------------------
//                                                                           ScheduledItem
//----------------------------------------------------------------------------------------


type ScheduledItem struct {
    Interval    time.Duration
    Task        func()
    nextTime    time.Time
}
type ScheduledItems []ScheduledItem


func (si ScheduledItems) Len() int {
    return len(si)
}

func (si ScheduledItems) Swap(i, j int) {
    t := si[i]
    si[i] = si[j]
    si[j] = t
}

func (si ScheduledItems) Less(i, j int) bool {
    return si[i].nextTime.Before(si[j].nextTime)
}

func (si *ScheduledItems) Push(x interface{}) {
    *si = append(*si, x.(ScheduledItem))
}

func (si *ScheduledItems) Pop() interface{} {
    old := *si
    n := len(old)
    x := old[n-1]
    *si = old[0 : n-1]
    return x
}

func (si ScheduledItems) Head() *ScheduledItem {
    if si.Len() > 0 {
        return &si[0]
    } else {
        return nil
    }
}


//----------------------------------------------------------------------------------------
//                                                                       Scheduler Control
//----------------------------------------------------------------------------------------


var schedulerChannel chan ScheduledItem


//  Two possible problems:
//  - Long running tasks can get re-scheduled before finished.
//  - Short-intervaled tasks can starve other
//    tasks by always being scheduled first.

func scheduler() {
    Log.LogFunctionName()
    defer Log.Debugf("=> Exit Scheduler <=")

    scheduledItems := new(ScheduledItems)
    heap.Init(scheduledItems)

    var shouldContinue bool = true
    for shouldContinue {

        var item *ScheduledItem
        if scheduledItems.Len() > 0 {
            item = scheduledItems.Head()
        }

        for item != nil && item.nextTime.Before(time.Now()) {
            go item.Task()
            item.nextTime = time.Now().Add(item.Interval)
            heap.Fix(scheduledItems, 0)
            item = scheduledItems.Head()
        }

        var waitTime time.Duration = time.Hour
        if  item != nil {
            waitTime = time.Since(item.nextTime) * -1
        }

        var timer *time.Timer = time.NewTimer(waitTime)
        select {

        case newItem := <- schedulerChannel:
            Log.Debugf("New scheduler item.")
            if newItem.Interval < 0 {
                shouldContinue = false
            } else {
                newItem.nextTime = time.Now().Add(newItem.Interval)
                heap.Push(scheduledItems, newItem)
            }

        case <- timer.C:
            //Log.Debugf("Run scheduler item.")
        }
    }
}


func StartScheduler() {
    Log.LogFunctionName()
    schedulerChannel = make(chan ScheduledItem)
    go scheduler()
}


func StopScheduler() {
    Log.LogFunctionName()
    item := ScheduledItem {
        Interval:   -1,
    }
    schedulerChannel <- item
}


func ScheduleTask(frequency time.Duration, task func()) {
    Log.LogFunctionName()
    if frequency <= 0 {
        frequency = time.Millisecond * 100
    }
    item := ScheduledItem {
        Interval:   frequency,
        Task:       task,
    }
    schedulerChannel <- item
}

