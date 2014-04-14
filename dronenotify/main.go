/*
Package dronenotify is an executeable that reads the environment variables
and sends an HTTP POST with a JSON payload to a cinotify server.

This variable is required for the bot to connect to an endpoint and must be set.
DRONE_NOTIFY_ADDRESS - eg localhost:3333

The variables it reads and json-ify's are as follows:
DRONE_REPO_SLUG
DRONE_BUILD_URL
DRONE_BUILD_DIR
DRONE_BUILD_NUMBER
DRONE_COMMIT
DRONE_BRANCH
*/
package main

import (
	"bytes"
	"encoding/json"
	"github.com/aarondl/cinotify/drone"
	"log"
	"net/http"
	"os"
)

func main() {
	addr := os.Getenv("DRONE_NOTIFY_ADDRESS")

	if len(addr) == 0 {
		log.Printf("DRONE_NOTIFY_ADDRESS not set, silently failing.")
		return
	}

	dr := drone.Notification{
		RepoSlug:    os.Getenv("DRONE_REPO_SLUG"),
		BuildUrl:    os.Getenv("DRONE_BUILD_URL"),
		BuildDir:    os.Getenv("DRONE_BUILD_DIR"),
		BuildNumber: os.Getenv("DRONE_BUILD_NUMBER"),
		Commit:      os.Getenv("DRONE_COMMIT"),
		Branch:      os.Getenv("DRONE_BRANCH"),
	}

	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	err := encoder.Encode(dr)

	if err != nil {
		log.Fatalf("dronenotify: Failed to encode json %v", err)
	}

	addr = "http://" + addr + "/"
	req, err := http.NewRequest("POST", addr, buffer)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "dronenotify")
	if err != nil {
		log.Fatalf("dronenotify: Failed to create request %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("dronenotify: Failed to send request %v", err)
	}

	log.Printf("dronenotify: Server responded with HTTP %v", resp.Status)
}
