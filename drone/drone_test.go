package drone

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aarondl/cinotify"
)

var testNotification = Notification{
	RepoSlug:    "repo_slug",
	BuildURL:    "build_url",
	BuildDir:    "build_dir",
	BuildNumber: "build_number",
	Commit:      "commit",
	Branch:      "branch",
}

func TestString(t *testing.T) {
	expect := "Drone[repo_slug]: Job #build_number Initiated at build_url"

	if got := testNotification.String(); got != expect {
		t.Errorf("Expected: %s, got: %s", expect, got)
	}
}

func TestHandle(t *testing.T) {
	var err error
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	if err = encoder.Encode(testNotification); err != nil {
		t.Error("Failed to jsonify payload: ", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/", buf)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "dronenotify")

	d := droneHandler{}
	note := d.Handle(req)
	if note == nil {
		t.Error("Expected to get a notification, got nil.")
	}

	if notification, ok := note.(Notification); !ok {
		t.Error("Expected to get a Notification type.")
	} else if notification != testNotification {
		t.Error("Expected an unaltered payload.")
	}
}

func TestHandleFail(t *testing.T) {
	buf := bytes.NewBufferString("{!$@($*&@&$)(*$)*&@$)")
	logger := &bytes.Buffer{}

	cinotify.Logger = log.New(logger, "", log.LstdFlags)

	req := httptest.NewRequest(http.MethodPost, "/", buf)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "dronenotify")

	if 0 != logger.Len() {
		t.Error("How could something be logged at this point?")
	}

	d := droneHandler{}
	note := d.Handle(req)
	if note != nil {
		t.Error("Expected an error to occur.")
	}

	if 0 == logger.Len() {
		t.Error("Expected something to be written to the log.")
	}
}
