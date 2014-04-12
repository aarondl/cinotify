package cinotify

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// DroneRequest is the fields transmitted from a dronenotify request.
type DroneRequest struct {
	RepoSlug    string `json:"repo_slug"`
	BuildUrl    string `json:"build_url"`
	BuildDir    string `json:"build_dir"`
	BuildNumber string `json:"build_number"`
	Commit      string `json:"commit"`
	Branch      string `json:"branch"`
}

// String converts a coverallsRequest to a tidy string for human consumption.
func (dr DroneRequest) String() string {
	return fmt.Sprintf(
		"Drone[%v]: Job #%v Initiated at %v",
		dr.RepoSlug,
		dr.BuildNumber,
		dr.BuildUrl,
	)
}

// droneHandler handles any requests from drone.
func droneHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	var dr DroneRequest
	err := decoder.Decode(&dr)
	if err != nil {
		doLogf("cinotify: Failed to decode json: %v", err)
		return
	}

	dispatchDroneCallbacks(&dr)
}
