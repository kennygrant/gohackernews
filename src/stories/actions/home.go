package storyactions

import (
	"fmt"
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

// HandleHome displays a list of stories using gravity to order them
// used for the home page for gravity rank see votes.go
// responds to GET /
func HandleHome(w http.ResponseWriter, r *http.Request) error {
	stats.RegisterHit(r)

	// Build a query
	q := stories.Query().Limit(listLimit)

	// Select only above 0 points,  Order by rank, then points, then name
	q.Where("points > 0").Order("rank desc, points desc, id desc")

	// Fetch the  params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Set the offset in pages if we have one
	page := int(params.GetInt("page"))
	if page > 0 {
		q.Offset(listLimit * page)
	}

	// Fetch the stories
	results, err := stories.FindAll(q)
	if err != nil {
		return server.InternalError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.AddKey("page", page)
	view.AddKey("stories", results)
	view.Template("stories/views/index.html.got")
	view.AddKey("pubdate", storiesModTime(results))
	view.AddKey("meta_title", fmt.Sprintf("%s, %s", config.Get("meta_title"), config.Get("meta_desc")))
	view.AddKey("meta_desc", config.Get("meta_desc"))
	view.AddKey("meta_keywords", config.Get("meta_keywords"))
	view.AddKey("meta_rss", storiesXMLPath(w, r))
	view.AddKey("userCount", stats.UserCount())
	view.AddKey("currentUser", session.CurrentUser(w, r))
	// For rss feeds use xml templates
	if strings.HasSuffix(r.URL.Path, ".xml") {
		view.Layout("")
		view.Template("stories/views/index.xml.got")
	}

	return view.Render()

}
