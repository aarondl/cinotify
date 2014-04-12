#dronenotify

Because drone.io has no webhook functionality you can use this package to
push notifications to your awaiting service.

In your drone.io configuration add:

```bash
go get "github.com/aarondl/cinotify/dronenotify"
```

This should give you access to dronenotify command. Drone.io sets almost
all the environment variables it needs (for details on which environment
variables are read, see the package documentation).

Please remember to set the following environment variable:

```
DRONE_NOTIFY_ADDRESS = yourserver.com:3333
```

Lastly in your drone.io configuration you can run the command:

```bash
dronenotify
```
