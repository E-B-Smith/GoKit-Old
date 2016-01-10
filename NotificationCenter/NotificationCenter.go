

package NotificationCenter


import (
    "fmt"
    "reflect"
    "violent.blue/GoKit/Log"
)


type NotificationFunction func(eventName string, eventData interface{})


type NotificationCenter struct {
    eventNameMapChannel     chan *map[string][]NotificationFunction
    eventNameMapChannelDone chan bool
    isRunning               bool
    isDone                  bool
}


func (center *NotificationCenter) eventNameMapServer() {
    eventNameMap := make(map[string][]NotificationFunction)
    isRunning := center.isRunning
    for isRunning {
        center.eventNameMapChannel <- &eventNameMap
        isRunning = <- center.eventNameMapChannelDone
    }
    center.isDone = true
}


func (center *NotificationCenter) acquireEventNameMap() *map[string][]NotificationFunction {
    if center.isDone { panic(fmt.Errorf("NotificationCenter already closed.")) }
    return <- center.eventNameMapChannel
}


func (center *NotificationCenter) releaseEventNameMap() {
    if center.isDone { return }
    center.eventNameMapChannelDone <- center.isRunning
}


func NewNotificationCenter() *NotificationCenter {
    var center NotificationCenter
    center.eventNameMapChannel = make(chan *map[string][]NotificationFunction)
    center.eventNameMapChannelDone = make(chan bool)
    center.isRunning = true
    go center.eventNameMapServer()
    return &center
}


func (center *NotificationCenter) Close() {
    Log.Debugf("NotificationCenter closed.")
    center.isRunning = false
    center.acquireEventNameMap()
    center.releaseEventNameMap()
}


func (center *NotificationCenter) PostNotification(eventName string, eventData interface{}) {
    eventNameMap := center.acquireEventNameMap()
    defer center.releaseEventNameMap()

    if eventObservers, ok := (*eventNameMap)[eventName]; ok {
        for _, observerFunc := range eventObservers {
            observerFunc(eventName, eventData)
        }
    }
}


func removeEventObserverGuts(eventNameMap *map[string][]NotificationFunction, eventName string, funIn NotificationFunction) {
    var ok bool
    var observerArray []NotificationFunction
    observerArray, ok = (*eventNameMap)[eventName]
    if !ok { return }

    pIn := reflect.ValueOf(funIn).Pointer()

    for idx, notificationFunc := range observerArray {
        if  reflect.ValueOf(notificationFunc).Pointer() == pIn {
            observerArray = append(observerArray[:idx], observerArray[idx+1:]...)
            (*eventNameMap)[eventName] = observerArray
            return
        }
    }
}


func (center *NotificationCenter) AddObserver(eventName string, funIn NotificationFunction) {
    if funIn == nil { return }

    eventNameMap := center.acquireEventNameMap()
    defer center.releaseEventNameMap()

    removeEventObserverGuts(eventNameMap, eventName, funIn)

    var ok bool
    var observerArray []NotificationFunction
    observerArray, ok = (*eventNameMap)[eventName]
    if !ok { observerArray = make([]NotificationFunction, 0) }

    observerArray = append(observerArray, funIn)
    (*eventNameMap)[eventName] = observerArray
}


func (center *NotificationCenter) RemoveEventObserver(eventName string, funIn NotificationFunction) {
    if funIn == nil { return }
    eventNameMap := center.acquireEventNameMap()
    defer center.releaseEventNameMap()
    removeEventObserverGuts(eventNameMap, eventName, funIn)
}


func (center *NotificationCenter) RemoveObserver(funIn NotificationFunction) {
    if funIn == nil { return }

    eventNameMap := center.acquireEventNameMap()
    defer center.releaseEventNameMap()

    for eventName, _ := range *eventNameMap {
        removeEventObserverGuts(eventNameMap, eventName, funIn)
    }
}

