/*
Package cinotify is a library that can listen for HTTP a webhook
from a cloud service. It then delivers those HTTP requests to the appropriate
extension for processing, and then on to any clients registered to deal with
those notifications.
*/
package cinotify

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// Handler is capable of handling a webhook from a given service.
type Handler interface {
	// Route is called before webserver starts add criteria to a route that will
	// be able to uniquely route this request to this notifier. Try and use
	// additional headers and user agent to differentiate against other
	// notifiers.
	//
	// Make sure you use at least route.Path("/") in your criteria.
	Route(*mux.Route)

	// Handle takes a request and transforms it into a Notifier.
	Handle(r *http.Request) fmt.Stringer
}

// handler holds a little bit of meta data about a Handler.
type handler struct {
	name        string
	realHandler Handler
	notifiers   []Notifier
	notifyfuncs []NotifyFunc
}

var handlers = make(map[string]*handler)

// Register registers an extension with cinotify.
func Register(name string, h Handler) {
	handlers[name] = &handler{
		name:        name,
		realHandler: h,
		notifiers:   make([]Notifier, 0),
		notifyfuncs: make([]NotifyFunc, 0),
	}
}

// Notifier takes a notification and notifies someone that something happened.
type Notifier interface {
	// Notify is called when a notification has been received. name will contain
	// the name of the extension that was called, and the notification will
	// be a fmt.Stringer you can type assert into whatever extension's custom
	// types.
	Notify(name string, notification fmt.Stringer)
}

// NotifyFunc is called when a notification has been received. name will contain
// the name of the extension that was called, and the notification will
// be a fmt.Stringer you can type assert into whatever extension's custom
// types.
type NotifyFunc func(name string, notification fmt.Stringer)

// To is used to direct notifications to a Notifier.
func To(n Notifier) {
	for _, h := range handlers {
		h.notifiers = append(h.notifiers, n)
	}
}

// ToFunc is used to direct notifications to a NotifyFunc.
func ToFunc(n NotifyFunc) {
	for _, h := range handlers {
		h.notifyfuncs = append(h.notifyfuncs, n)
	}
}

// context is a dumb helper object for a useless fluent syntax I'm creating,
// see When().
type context struct {
	name string
}

// To is used to direct notifications from the context's service to a Notifier.
func (c context) To(n Notifier) {
	for name, h := range handlers {
		if name == c.name {
			h.notifiers = append(h.notifiers, n)
		}
	}
}

// ToFunc is used to direct notifications from the context's service to a
// NotifyFunc.
func (c context) ToFunc(n NotifyFunc) {
	for name, h := range handlers {
		if name == c.name {
			h.notifyfuncs = append(h.notifyfuncs, n)
		}
	}
}

// When is a matcher for an extension name, use it in front of a to call to
// limit which notifications your Notifier/NotifyFunc will get.
// Example: When("github").To(myHandler)
func When(name string) context {
	return context{name}
}

// createRouter takes all the registered handlers and creates proxy anon
// functions for them, routing all the requests to the correct spot.
func createRouter() *mux.Router {
	router := mux.NewRouter()

	for name, h := range handlers {
		r := router.NewRoute()
		h.realHandler.Route(r)
		r.Name(name)
		r.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			notification := h.realHandler.Handle(r)
			dispatch(notification)
		})
	}

	return router
}

// dispatch sends the notifications out to all appropriate Notifier and
// NotifyFuncs.
func dispatch(notification fmt.Stringer) {
	for name, h := range handlers {
		for _, notifier := range h.notifiers {
			notifier.Notify(name, notification)
		}
		for _, notifyfunc := range h.notifyfuncs {
			notifyfunc(name, notification)
		}
	}
}

// StartServer starts listening on the given port, if 0 will default to 5000.
func StartServer(port uint16) <-chan error {
	address := ":" + strconv.Itoa(int(port))

	router := createRouter()

	ch := make(chan error, 1)
	go func() {
		DoLogf("listening on: [%v]", address)
		err := http.ListenAndServe(address, router)
		DoLogf("cinotify: error listening %v", err)
		ch <- err
	}()
	return ch
}
