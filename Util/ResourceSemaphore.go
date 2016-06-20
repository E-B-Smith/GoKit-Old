

//----------------------------------------------------------------------------------------
//
//                                                                    ResourceSemaphore.go
//                                           A semaphore for limiting resource utilization
//
//                                                                    E.B.Smith, June 2016
//                        -©- Copyright © 2014-2016 Edward Smith, all rights reserved. -©-
//
//----------------------------------------------------------------------------------------


package Util


import (
    "sync"
)


/*

    Example:

    resourceLimit := NewResourceSemphoreWithLimit(20)
    for resourceLimit.ActiveCount() < resourceLimit.TotalCount() && ! done {

        resourceLimit.BeginAddToSemaphore()
        go func(page) {
            parsePage(page)
            resourceLimit.Done()
        } (page)
        resourceLimit.EndAddToSemaphore()

        resourceLimit.WaitAvailable()
    }
    resourceLimit.WaitComplete()

*/


type ResourceSemaphore struct {
    threadLock  sync.Mutex
    threadSync  *sync.Cond
    threadCount int
    threadMax   int
}


func (r *ResourceSemaphore) ActiveCount() int {
    r.threadLock.Lock()
    defer r.threadLock.Unlock()
    return r.threadCount
}


func (r *ResourceSemaphore) TotalCount() int {
    return r.threadMax
}


func (r *ResourceSemaphore) SetTotal(total int) {
    r.threadMax = total
}


func (r *ResourceSemaphore) BeginAddToSemaphore() {
    r.threadLock.Lock()
    r.threadCount++
}


func (r *ResourceSemaphore) EndAddToSemaphore() {
    r.threadLock.Unlock()
}


func (r *ResourceSemaphore) WaitAvailable() {
    r.threadLock.Lock()
    for r.threadCount >= r.threadMax {
        r.threadSync.Wait()     //  Wait unlocks then locks.
    }
    r.threadLock.Unlock()
}


func (r *ResourceSemaphore) WaitComplete() {
    r.threadLock.Lock()
    for r.threadCount > 0 {
        r.threadSync.Wait()     //  Wait unlocks then locks.
    }
    r.threadLock.Unlock()
}


func (r *ResourceSemaphore) Done() {
    r.threadLock.Lock()
    r.threadCount--
    r.threadLock.Unlock()
    r.threadSync.Broadcast()
}


func NewResourceSemphore() *ResourceSemaphore {
    var r ResourceSemaphore
    r.threadSync = sync.NewCond(&r.threadLock)
    return &r;
}


func NewResourceSemphoreWithLimit(limit int) *ResourceSemaphore {
    r := NewResourceSemphore()
    r.threadMax = limit
    return r
}

