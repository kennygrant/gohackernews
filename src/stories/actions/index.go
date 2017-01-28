package storyactions

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fragmenta/mux"
	"github.com/fragmenta/server"
	"github.com/fragmenta/server/config"
	"github.com/fragmenta/view"

	"github.com/kennygrant/gohackernews/src/lib/session"
	"github.com/kennygrant/gohackernews/src/lib/stats"
	"github.com/kennygrant/gohackernews/src/stories"
)

const listLimit = 50

// storiesModTime returns the mod time of the first story, or current time if no stories
func storiesModTime(availableStories []*stories.Story) time.Time {
	if len(availableStories) == 0 {
		return time.Now()
	}
	story := availableStories[0]

	return story.UpdatedAt
}

// storiesXMLPath returns the xml path for a given request to a stories link
func storiesXMLPath(w http.ResponseWriter, r *http.Request) string {

	p := strings.Replace(r.URL.Path, ".xml", "", 1)
	if p == "/" {
		p = "/index"
	}

	q := r.URL.RawQuery
	if len(q) > 0 {
		q = "?" + q
	}

	return fmt.Sprintf("%s.xml%s", p, q)
}

// HandleIndex displays a list of stories.
func HandleIndex(w http.ResponseWriter, r *http.Request) error {

	// Authorise list story - anyone can view stories

	// Get the params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	stats.RegisterHit(r)

	// Build a query
	q := stories.Query().Limit(listLimit)

	// Order by date by default
	q.Where("points > -6").Order("created_at desc")

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

	windowTitle := config.Get("meta_title")
	switch filter {
	case "Video:":
		windowTitle = "Golang Videos"
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.AddKey("page", page)
	view.AddKey("stories", results)
	view.AddKey("pubdate", storiesModTime(results))
	view.AddKey("meta_title", windowTitle)
	view.AddKey("meta_desc", config.Get("meta_desc"))
	view.AddKey("meta_keywords", config.Get("meta_keywords"))
	view.AddKey("meta_rss", storiesXMLPath(w, r))
	view.AddKey("currentUser", session.CurrentUser(w, r))

	if strings.HasSuffix(r.URL.Path, ".xml") {
		view.Layout("")
		view.Template("stories/views/index.xml.got")
	}

	return view.Render()

}
