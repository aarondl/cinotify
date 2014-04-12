package cinotify

import (
	"fmt"
	"github.com/gorilla/schema"
	"net/http"
)

// decoder helps decode all the post form requests
var decoder = schema.NewDecoder()

// CoverallsRequest is the fields passed back from the coveralls webhook.
type CoverallsRequest struct {
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

// String converts a coverallsRequest to a tidy string for human consumption.
func (cr CoverallsRequest) String() string {
	change := ""
	if cr.CoverageChange == 1.0 {
		change += "none"
	} else if cr.CoverageChange > 1.0 {
		change = fmt.Sprintf("+%.0f%%", (cr.CoverageChange-1.0)*100.0)
	} else {
		change = fmt.Sprintf("-%v%%", (1.0-cr.CoverageChange)*100.0)
	}

	return fmt.Sprintf(
		"Coveralls[%s]: Change(%s) Percent(%.2f) %s",
		cr.RepoName,
		change,
		cr.CoveredPercent,
		cr.Url,
	)
}

// coverallsHandler handles any requests from coveralls.
func coverallsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	err := r.ParseForm()
	if err != nil {
		doLogf("cinotify: Failed to parse form body: %v", err)
		return
	}

	var cr CoverallsRequest
	err = decoder.Decode(&cr, r.PostForm)
	if err != nil {
		doLogf("cinotify: Failed to decode post form: %v", err)
		return
	}

	dispatchCoverallsCallbacks(&cr)
}
