package storyactions

import (
	"net/http"
	"strings"

	"github.com/fragmenta/mux"
	"github.com/fragmenta/server"
	"github.com/fragmenta/server/config"
	"github.com/fragmenta/view"

	"github.com/kennygrant/gohackernews/src/lib/session"
	"github.com/kennygrant/gohackernews/src/lib/stats"
	"github.com/kennygrant/gohackernews/src/stories"
)

// HandleJobs repsponds to GET /go-jobs
func HandleJobs(w http.ResponseWriter, r *http.Request) error {

	// No Authorisation - anyone can view stories

	// Get the params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	stats.RegisterHit(r)

	// Build a query
	q := stories.Query().Limit(listLimit)

	// Filter for points
	q.Where("points > 0")

	// If filtering, order by rank, not by date
	q.Order("rank desc, points desc, created_at desc")

	// Filter on hiring title
	q.Where("stories.name LIKE 'Hiring:%'")

	// Filter if necessary - this assumes name and summary cols
	filter := params.Get("q")
	if len(filter) > 0 {

		// Replace special characters with escaped sequence
		filter = strings.Replace(filter, "_", "\\_", -1)
		filter = strings.Replace(filter, "%", "\\%", -1)

		wildcard := "%" + filter + "%"

		// Perform a wildcard search for name or url
		q.Where("stories.name ILIKE ? OR stories.url ILIKE ?", wildcard, wildcard)

		// If filtering, order by rank, not by date
		q.Order("rank desc, points desc, id desc")
	}

	// Set the offset in pages if we have one
	page := params.GetInt("page")
	if page > 0 {
		q.Offset(listLimit * int(page))
	}

	// Fetch the stories
	results, err := stories.FindAll(q)
	if err != nil {
		return server.InternalError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.Template("stories/views/index.html.got")
	view.AddKey("page", page)
	view.AddKey("stories", results)
	view.AddKey("pubdate", storiesModTime(results))
	view.AddKey("meta_title", "Go Jobs - Companies hiring programmers and using Go")
	view.AddKey("meta_desc", "Jobs for Go hackers")
	view.AddKey("meta_keywords", "jobs careers "+config.Get("meta_keywords"))
	view.AddKey("meta_rss", storiesXMLPath(w, r))
	view.AddKey("currentUser", session.CurrentUser(w, r))

	if strings.HasSuffix(r.URL.Path, ".xml") {
		view.Layout("")
		view.Template("stories/views/index.xml.got")
	}

	return view.Render()

}
