#cinotify

cinotify is a package that enables ci environments in the cloud to notify
a listener of events. Although it could be any service that supports
web hooks.

A good example (see code below) is if you wanted to know whenever drone.io
kicked off a build. You could display the link, build number etc.

Note that it's built extensibly similar to the sql/db package in the Go standard
library. That's to say: to enable an extension, import it with the _ identifier
(unless you plan on using the notification struct directly.)

###Extensions

| Extension    | Import Path |
| ------------ | ----------- |
| drone.io     | [github.com/aarondl/cinotify/drone](https://github.com/aarondl/cinotify/drone) |
| coveralls.io | [github.com/aarondl/cinotify/coveralls](https://github.com/aarondl/cinotify/coveralls) |

###Installation

```bash
go get github.com/aarondl/cinotify
```

###Usage

For example a program to listen for drone.io notifications could look like:

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aarondl/cinotify"
	_ "github.com/aarondl/cinotify/drone"
)

func main() {
	// Set logger
	cinotify.Logger = log.New(os.Stdout, "", log.LstdFlags)

	// Add any callbacks we need.
	cinotify.ToFunc(func(name string, notification fmt.Stringer) {
		log.Println(notification)
		// OR remove the _ from in front of drone's import and do:
		if droneNotification, ok := notification.(drone.Notification); ok {
			// Here we can access all the fields of the drone.Request struct.
		}
	})

	// Start server.
	ch := cinotify.StartServer(5000)
	log.Println(<-ch)
}
```

###Other Usages

To register callbacks there is 4 methods:

| Method | Use |
| ------ | --- |
| [To](http://godoc.org/github.com/aarondl/cinotify#To) | Use with  [Notifiers](http://godoc.org/github.com/aarondl/cinotify#Notifier). |
| [ToFunc](http://godoc.org/github.com/aarondl/cinotify#To) | Use with  [NotifyFunc](http://godoc.org/github.com/aarondl/cinotify#NotifyFunc) functions. |
| [When(name).To](http://godoc.org/github.com/aarondl/cinotify#When) | Same as To above, but restricts events to the service name, eg: [drone.Name](http://godoc.org/github.com/aarondl/cinotify/drone#Name) |
| [When(name).ToFunc](http://godoc.org/github.com/aarondl/cinotify#When) | Same as When().To above but uses functions. |
