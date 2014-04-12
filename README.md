#cinotify

cinotify is a package that enables ci environments in the cloud to notify
a listener of events.

A good example (see code below) is if you wanted to know whenever drone.io
kicked off a build. You could display the link, build number etc.

__Supported Services:__ drone.io coveralls.io

###Installation

```bash
go get github.com/aarondl/cinotify
```

###Usage

For example a simple server to listen for drone.io requests could look like:

```go
package main

import "github.com/aarondl/cinotify"
import "log"
import "os"

func main() {
	// Set logger
	cinotify.Logger = log.New(os.Stdout, "", log.LstdFlags)

	// Add any callbacks we need.
	cinotify.AddDroneCallback(func(dr *cinotify.DroneRequest) {
		log.Println(dr)
	})

	// Start server.
	ch := cinotify.StartServer(5000, cinotify.Drone)

	log.Println(<-ch)
}
```
