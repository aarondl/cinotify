package cinotify

import (
	"net/http"
	"time"
)

// StartServer starts listening on the address.
//
// This blocks so you should start it in a goroutine if you want
func StartServer(address string) error {
	server := http.Server{
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		IdleTimeout:  time.Second * 5,

		Addr:    address,
		Handler: http.HandlerFunc(serveHTTP),
	}

	DoLogf("listening on: [%v]", address)
	err := server.ListenAndServe()
	if err != nil {
		DoLogf("cinotify: error listening %v", err)
	}

	return err
}

func serveHTTP(w http.ResponseWriter, r *http.Request) {
	foundOne := false

	for _, h := range handlers {
		stringer := h.realHandler.Handle(r)
		if stringer == nil {
			continue
		}

		foundOne = true

		for _, n := range h.notifiers {
			n.Notify(h.name, stringer)
		}
	}

	if !foundOne {
		DoLogf("cinotify: failed to find handler for request (HTTP %s %s):\n%v\n",
			r.Method,
			r.URL.Path,
			r.Header,
		)
	}
}
