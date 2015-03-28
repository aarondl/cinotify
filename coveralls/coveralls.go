// Package coveralls is an extension for the cinotify package. See
// github.com/aarondl/cinotify for details.
package coveralls

import (
	"fmt"
	"net/http"

	"github.com/aarondl/cinotify"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

// Name is the name of the service, for use with When() in cinotify.
const Name = "coveralls"

func init() {
	cinotify.Register(Name, coverallsHandler{})
}

// decoder helps decode all the post form requests
var decoder = schema.NewDecoder()

// Notification is the notification that arrives from a coveralls webhook.
type Notification struct {
	BadgeUrl       string  `schema:"badge_url"`
	Branch         string  `schema:"branch"`
	CommitMessage  string  `schema:"commit_message"`
	CommitSha      string  `schema:"commit_sha"`
	CommitterEmail string  `schema:"committer_email"`
	CommitterName  string  `schema:"committer_name"`
	CoverageChange float64 `schema:"coverage_change"`
	CoveredPercent float64 `schema:"covered_percent"`
	RepoName       string  `schema:"repo_name"`
	Url            string  `schema:"url"`
}

// String converts a Notification to a tidy string for human consumption.
func (n Notification) String() string {
	return fmt.Sprintf(
		"Coveralls[%s]: Change(%.2f%%) Covered(%.2f%%) %s",
		n.RepoName,
		n.CoverageChange,
		n.CoveredPercent,
		n.Url,
	)
}

// coverallsHandler implements cinotify.Handler
type coverallsHandler struct {
}

// coverallsHandler handles any requests from coveralls.
func (_ coverallsHandler) Handle(r *http.Request) fmt.Stringer {
	defer r.Body.Close()

	err := r.ParseForm()
	if err != nil {
		cinotify.DoLogf("cinotify/coveralls: Failed to parse form: %v", err)
		return nil
	}

	var n Notification
	err = decoder.Decode(&n, r.PostForm)
	if err != nil {
		cinotify.DoLogf("cinotify/coveralls: Failed to decode form: %v", err)
		return nil
	}

	return n
}

// Route creates a route that only a coveralls client should hit.
func (_ coverallsHandler) Route(r *mux.Route) {
	r.Path("/").Methods("POST").Headers(
		"Content-Type", "application/x-www-form-urlencoded",
		"User-Agent", "Ruby",
	)
}
