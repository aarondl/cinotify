package cinotify

import (
	"fmt"
	"net/http"
)

// Global list of handlers
var (
	handlers = make(map[string]*handler)
)

// Handler is capable of handling a webhook from a given service.
type Handler interface {
	// Handle takes an http request and transform it into a fmt.Stringer.
	Handle(r *http.Request) fmt.Stringer
}

// HandlerFunc is the function version of Handler
type HandlerFunc func(r *http.Request) fmt.Stringer

// Handle delegates to the HandlerFunc's fn
func (h HandlerFunc) Handle(r *http.Request) fmt.Stringer {
	return h(r)
}

// handler holds a little bit of meta data about a Handler.
type handler struct {
	name        string
	realHandler Handler
	notifiers   []Notifier
}

// Register registers an extension with cinotify.
func Register(name string, h Handler) {
	handlers[name] = &handler{
		name:        name,
		realHandler: h,
		notifiers:   make([]Notifier, 0),
	}
}
