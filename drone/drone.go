// Package drone is an extension for the cinotify package. See
// github.com/aarondl/cinotify for details.
package drone

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aarondl/cinotify"
)

// Name is the name of the service, for use with When() in cinotify.
const Name = "drone"

func init() {
	cinotify.Register(Name, droneHandler{})
}

// Notification is the notification transmitted from a dronenotify request.
type Notification struct {
	RepoSlug    string `json:"repo_slug"`
	BuildURL    string `json:"build_url"`
	BuildDir    string `json:"build_dir"`
	BuildNumber string `json:"build_number"`
	Commit      string `json:"commit"`
	Branch      string `json:"branch"`
}

// String converts a coverallsRequest to a tidy string for human consumption.
func (n Notification) String() string {
	return fmt.Sprintf(
		"Drone[%v]: Job #%v Initiated at %v",
		n.RepoSlug,
		n.BuildNumber,
		n.BuildURL,
	)
}

// droneHandler implements cinotify.Handler
type droneHandler struct {
}

// droneHandler handles any requests from drone.
func (droneHandler) Handle(r *http.Request) fmt.Stringer {
	if r.URL.Path != "/" || r.Method != http.MethodPost {
		return nil
	}
	if r.Header.Get("Content-Type") != "application/json" {
		return nil
	}
	if r.Header.Get("User-Agent") != "dronenotify" {
		return nil
	}

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	var n Notification
	err := decoder.Decode(&n)
	if err != nil {
		cinotify.DoLog("cinotify/drone: Failed to decode json payload: ", err)
		return nil
	}

	return n
}
