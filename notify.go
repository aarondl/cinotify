/*
Package cinotify is a library that can listen for an HTTP webhook
from a cloud service. It then delivers those HTTP requests to the appropriate
extension for processing, and then on to any clients registered to deal with
those notifications.

A quick usage example would be:

	import (
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
				// Here we can access all the fields of the drone.Notification struct
			}
		})

		// Start server.
		ch := cinotify.StartServer(5000)
		log.Println(<-ch)
	}
*/
package cinotify

import (
	"fmt"
)

// Notifier takes a notification and notifies someone that something happened.
type Notifier interface {
	// Notify is called when a notification has been received. name will contain
	// the name of the extension that was called, and the notification will
	// be a fmt.Stringer you can type assert into whatever extension's custom
	// types.
	Notify(name string, notification fmt.Stringer)
}

// NotifyFunc is the function version of the Notifier interface, see docs
// for Notifier.Notify.
type NotifyFunc func(name string, notification fmt.Stringer)

// Notify for a NotifyFunc
func (fn NotifyFunc) Notify(name string, notification fmt.Stringer) {
	fn(name, notification)
}

// To is used to direct notifications to a Notifier.
func To(n Notifier) {
	for _, h := range handlers {
		h.notifiers = append(h.notifiers, n)
	}
}

// Filter is a dumb helper object for a useless fluent syntax, see When().
type Filter struct {
	name string
}

// To is used to direct notifications from the filter's named service to a
// Notifier.
func (f Filter) To(n Notifier) {
	for name, h := range handlers {
		if name == f.name {
			h.notifiers = append(h.notifiers, n)
		}
	}
}

// When is a matcher for an extension name, use it in front of a To() call to
// limit which notifications your Notifier/NotifyFunc will get.
// Example: When(drone.Name).To(myHandler)
func When(name string) Filter {
	return Filter{name}
}

// dispatch sends the notifications out to all appropriate Notifier and
// NotifyFuncs.
func dispatch(name string, notification fmt.Stringer) {
	h := handlers[name]
	for _, notifier := range h.notifiers {
		notifier.Notify(name, notification)
	}
}
