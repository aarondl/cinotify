// Package coveralls is an extension for the cinotify package. See
// github.com/aarondl/cinotify for details.
package coveralls

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/aarondl/cinotify"
)

// Name is the name of the service, for use with When() in cinotify.
const Name = "coveralls"

func init() {
	cinotify.Register(Name, coverallsHandler{})
}

// Notification is the notification that arrives from a coveralls webhook.
type Notification struct {
	BadgeURL       string  `schema:"badge_url"`
	Branch         string  `schema:"branch"`
	CommitMessage  string  `schema:"commit_message"`
	CommitSha      string  `schema:"commit_sha"`
	CommitterEmail string  `schema:"committer_email"`
	CommitterName  string  `schema:"committer_name"`
	CoverageChange float64 `schema:"coverage_change"`
	CoveredPercent float64 `schema:"covered_percent"`
	RepoName       string  `schema:"repo_name"`
	URL            string  `schema:"url"`
}

// String converts a Notification to a tidy string for human consumption.
func (n Notification) String() string {
	return fmt.Sprintf(
		"Coveralls[%s]: Change(%.2f%%) Covered(%.2f%%) %s",
		n.RepoName,
		n.CoverageChange,
		n.CoveredPercent,
		n.URL,
	)
}

// coverallsHandler implements cinotify.Handler
type coverallsHandler struct {
}

// coverallsHandler handles any requests from coveralls.
func (coverallsHandler) Handle(r *http.Request) fmt.Stringer {
	if r.URL.Path != "/" || r.Method != http.MethodPost {
		return nil
	}
	if r.Header.Get("User-Agent") != "ruby" || r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		return nil
	}

	err := r.ParseForm()
	if err != nil {
		cinotify.DoLogf("cinotify/coveralls: Failed to parse form: %v", err)
		return nil
	}

	n := Notification{
		BadgeURL:       r.PostFormValue("badge_url"),
		Branch:         r.PostFormValue("branch"),
		CommitMessage:  r.PostFormValue("commit_message"),
		CommitSha:      r.PostFormValue("commit_sha"),
		CommitterEmail: r.PostFormValue("committer_email"),
		CommitterName:  r.PostFormValue("committer_name"),
		RepoName:       r.PostFormValue("repo_name"),
		URL:            r.PostFormValue("url"),
	}
	coverageChange := r.PostFormValue("coverage_change")
	if n.CoverageChange, err = strconv.ParseFloat(coverageChange, 10); err != nil {
		cinotify.DoLogf("cinotify/coveralls: Failed to parse coverage_change(%v): %q", err, coverageChange)
		return nil
	}
	coveredPercent := r.PostFormValue("covered_percent")
	if n.CoveredPercent, err = strconv.ParseFloat(coveredPercent, 10); err != nil {
		cinotify.DoLogf("cinotify/coveralls: Failed to parse coverage_percent(%v): %q", err, coveredPercent)
		return nil
	}

	return n
}
