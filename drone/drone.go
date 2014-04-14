package drone

import (
	"encoding/json"
	"fmt"
	"github.com/aarondl/cinotify"
	"github.com/gorilla/mux"
	"net/http"
)

// Name is the name of the service, for use with the When() in cinotify.
const Name = "drone"

func init() {
	cinotify.Register(Name, droneHandler{})
}

// Notification is the fields transmitted from a dronenotify request.
type Notification struct {
	RepoSlug    string `json:"repo_slug"`
	BuildUrl    string `json:"build_url"`
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
		n.BuildUrl,
	)
}

// droneHandler implements cinotify.Handler
type droneHandler struct {
}

// droneHandler handles any requests from drone.
func (_ droneHandler) Handle(r *http.Request) fmt.Stringer {
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

// Route creates a route that only a dronenotify client should hit.
func (_ droneHandler) Route(r *mux.Route) {
	r.Path("/").Methods("POST").Headers(
		"Content-Type", "application/json",
		"User-Agent", "dronenotify",
	)
}
