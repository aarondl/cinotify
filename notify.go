/*
Package cinotify is a library that can listen for HTTP requests (typically a
webhook) from a cloud service. It then delivers those requests to the
callbacks that were registered.

A quickstart might look like:

cinotify.Logger = someLogger
AddDroneCallback()
AddCoverallsCallback()
StartServer()
*/
package cinotify

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sync"
	"strconv"
)

// Service represents a CI service.
type Service string

// All the supported services.
const (
	Drone     Service = "drone.io"
	Coveralls Service = "coveralls.io"
)

// Logger variable can be set to have logged output for error handling / status.
// Logger must be set before StartServer is called.
var Logger *log.Logger

// callbackControl controlsl access to the lists of callbacks.
var callbackControl sync.RWMutex

// Callback types
type (
	// DroneCallback can be called when a Drone request is recieved.
	DroneCallback func(dr *DroneRequest)
	// CoverallsCallback can be called when a Coveralls request is recieved.
	CoverallsCallback func(dr *CoverallsRequest)
)

// The lists of callbacks
var (
	droneCallbacks = make([]DroneCallback, 0)
	coverallsCallbacks = make([]CoverallsCallback, 0)
)

// doLog logs anything of interest.
func doLog(args ...interface{}) {
	if Logger != nil {
		Logger.Print(args...)
	}
}

// doLogf logs anything of interest with a format
func doLogf(format string, args ...interface{}) {
	if Logger != nil {
		Logger.Printf(format, args...)
	}
}

// AddDroneCallback adds a callback to be called when drone.io recieves a
// request.
func AddDroneCallback(cb DroneCallback) {
	callbackControl.Lock()
	defer callbackControl.Unlock()
	droneCallbacks = append(droneCallbacks, cb)
}

// AddCoverallsCallback adds a callback to be called when coveralls.io recieves a
// request.
func AddCoverallsCallback(cb CoverallsCallback) {
	callbackControl.Lock()
	defer callbackControl.Unlock()
	coverallsCallbacks = append(coverallsCallbacks, cb)
}

// dispatchDroneCallbacks sends the DroneRequest to each callback.
func dispatchDroneCallbacks(dr *DroneRequest) {
	callbackControl.RLock()
	defer callbackControl.RUnlock()

	for _, fn := range droneCallbacks {
		fn(dr)
	}
}

// dispatchCoverallsCallbacks sends the CoverallsRequest to each callback.
func dispatchCoverallsCallbacks(cr *CoverallsRequest) {
	callbackControl.RLock()
	defer callbackControl.RUnlock()

	for _, fn := range coverallsCallbacks {
		fn(cr)
	}
}

// StartServer starts listening on the given port, if 0 will default to 5000.
// You must pass in a list of Services to enable, or the server will do nothing.
func StartServer(port uint16, enabled ...Service) <-chan error {
	address := ":" + strconv.Itoa(int(port))

	r := mux.NewRouter()

	for _, toEnable := range enabled {
		switch toEnable {
		case Drone:
			r.HandleFunc("/", droneHandler).
				Methods("POST").
				Headers("Content-Type", "application/json")
		case Coveralls:
			r.HandleFunc("/", coverallsHandler).
				Methods("POST").
				Headers("Content-Type", "application/x-www-form-urlencoded")
		}
	}

	ch := make(chan error, 1)
	go func() {
		doLogf("listening on: [%v]", address)
		err := http.ListenAndServe(address, r)
		doLogf("cinotify: error listening %v", err)
		ch <- err
	}()
	return ch
}
