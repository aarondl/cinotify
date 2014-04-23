package coveralls

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	. "testing"

	"github.com/aarondl/cinotify"
	"github.com/gorilla/mux"
)

var testNotification = Notification{
	BadgeUrl:       "badge_url",
	Branch:         "branch",
	CommitMessage:  "commit_message",
	CommitSha:      "commit_sha",
	CommitterEmail: "committer_email",
	CommitterName:  "committer_name",
	CoverageChange: 1.5,
	CoveredPercent: 97.0,
	RepoName:       "repo_name",
	Url:            "url",
}

func TestString(t *T) {
	expect := "Coveralls[repo_name]: Change(1.50%) Covered(97.00%) url"

	if got := testNotification.String(); got != expect {
		t.Error("Expected: %s, got: %s", expect, got)
	}
}

func TestHandle(t *T) {
	var err error
	vals := url.Values{}
	vals.Add("badge_url", "badge_url")
	vals.Add("branch", "branch")
	vals.Add("commit_message", "commit_message")
	vals.Add("commit_sha", "commit_sha")
	vals.Add("committer_email", "committer_email")
	vals.Add("committer_name", "committer_name")
	vals.Add("coverage_change", "1.5")
	vals.Add("covered_percent", "97.0")
	vals.Add("repo_name", "repo_name")
	vals.Add("url", "url")

	buf := bytes.NewBufferString(vals.Encode())

	var req *http.Request
	if req, err = http.NewRequest("POST", "/", buf); err != nil {
		t.Error("Error creating mock request: ", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	d := coverallsHandler{}
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

func TestHandleFail(t *T) {
	var err error
	buf := bytes.NewBufferString("{!$@($*&@&$)(*$)*&@$)")
	logger := &bytes.Buffer{}

	cinotify.Logger = log.New(logger, "", log.LstdFlags)

	var req *http.Request
	if req, err = http.NewRequest("POST", "/", buf); err != nil {
		t.Error("Error creating mock request: ", err)
	}

	if 0 != logger.Len() {
		t.Error("How could something be logged at this point?")
	}

	d := coverallsHandler{}
	note := d.Handle(req)
	if note != nil {
		t.Error("Expected an error to occur.")
	}

	if 0 == logger.Len() {
		t.Error("Expected something to be written to the log.")
	}
}

func TestRoute(t *T) {
	var err error

	d := coverallsHandler{}
	router := mux.NewRouter()
	r := router.NewRoute()

	d.Route(r)
	r.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	resp := httptest.NewRecorder()
	var req *http.Request
	if req, err = http.NewRequest("POST", "/", nil); err != nil {
		t.Error("Error creating mock request: ", err)
	}
	req.Header.Add("User-Agent", "Ruby")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	router.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Error("Route did not match request.")
	}
}
