//----------------------------------------------------------------------------------------
//
//                                                               ResourceSemaphore_Test.go
//                                           A semaphore for limiting resource utilization
//
//                                                                    E.B.Smith, June 2016
//                        -©- Copyright © 2014-2016 Edward Smith, all rights reserved. -©-
//
//----------------------------------------------------------------------------------------


package Util


import (
    "fmt"
    "time"
    "testing"
)


//----------------------------------------------------------------------------------------
//                                                                   TestResourceSemaphore
//----------------------------------------------------------------------------------------


const kMaxCount = 20


func TestResourceSemaphore(t *testing.T) {

    waiter := func() {
        time.Sleep(time.Millisecond*10)
    }

    assertActiveCount := func(activeCount, minCount, maxCount int) {
        if activeCount < minCount || activeCount > maxCount {
            t.Errorf("Expected count to be %d-%d. Got %d.\n", minCount, maxCount, activeCount)
        }
    }

    var i int
    resourceLimit := NewResourceSemphoreWithLimit(kMaxCount)
    for i := 0; i < kMaxCount; i++ {
        resourceLimit.WaitAvailable()
        resourceLimit.BeginAddToSemaphore()
        go func() {
            waiter()
            resourceLimit.Done()
        } ()
        resourceLimit.EndAddToSemaphore()
    }

    assertActiveCount(resourceLimit.ActiveCount(), kMaxCount, kMaxCount)

    i = 0
    for i < 10000 {

        i++
        assertActiveCount(resourceLimit.ActiveCount(), 0, kMaxCount)

        resourceLimit.WaitAvailable()
        assertActiveCount(resourceLimit.ActiveCount(), 0, kMaxCount-1)
        resourceLimit.BeginAddToSemaphore()
        go func() {
            waiter()
            resourceLimit.Done()
        } ()
        resourceLimit.EndAddToSemaphore()
    }

    fmt.Printf("Waiting for %d threads.\n", resourceLimit.ActiveCount())
    assertActiveCount(resourceLimit.ActiveCount(), kMaxCount, kMaxCount)
    resourceLimit.WaitComplete()
    assertActiveCount(resourceLimit.ActiveCount(), 0, 0)
    fmt.Printf("Waiting for %d threads.\n", resourceLimit.ActiveCount())

}

