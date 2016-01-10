

package NotificationCenter


import (
    "fmt"
    "time"
    "testing"
)


func TestNotifications1(t *testing.T) {

    callback := func(event string, data interface{}) {
        fmt.Printf("Event '%s' happened with data '%s'.\n", event, data)
    }

    center := NewNotificationCenter()
    //defer center.Close()
    fmt.Printf("Done: %t.\n", center.isDone)

    center.PostNotification("Happened", "Fail")

    center.AddObserver("Happened", callback)
    center.PostNotification("Happened", "Success")

    center.RemoveEventObserver("Happened", callback)
    center.PostNotification("Happened", "Fail")

    center.AddObserver("Happened", callback)
    center.PostNotification("Happened", "Success")

    center.RemoveObserver(callback)
    center.PostNotification("Happened", "Fail")

    center.AddObserver("Happened",
        func(eventName string, eventData interface{}) {
            fmt.Printf("Success!\n")
        });
    center.PostNotification("Happened", "Success")

    center.Close()
    time.Sleep(time.Second * 2.0)
    if ! center.isDone {
        t.Errorf("Not done!")
    }
}

