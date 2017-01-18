package storyactions

import (
	"net/http"
	"strings"

	"github.com/fragmenta/mux"
	"github.com/fragmenta/server"
	"github.com/fragmenta/server/config"
	"github.com/fragmenta/view"

	"github.com/kennygrant/gohackernews/src/lib/session"
	"github.com/kennygrant/gohackernews/src/stories"
)

// HandleListCode displays a list of stories linking to repos (github etc)
// responds to GET /stories/code
func HandleListCode(w http.ResponseWriter, r *http.Request) error {

	// Get params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Build a query
	q := stories.Query().Where("points > -6").Order("rank desc, points desc, id desc").Limit(listLimit)

	// Restrict to stories with have a url starting with github.com or bitbucket.org
	// other code repos can be added later
	q.Where("url ILIKE 'https://github.com%'").OrWhere("url ILIKE 'https://bitbucket.org%'")

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
	view.AddKey("pubdate", storiesModTime(results))
	view.AddKey("meta_title", "Go Code")
	view.AddKey("meta_desc", config.Get("meta_desc"))
	view.AddKey("meta_keywords", config.Get("meta_keywords"))
	view.AddKey("meta_rss", storiesXMLPath(w, r))
	view.Template("stories/views/index.html.got")
	view.AddKey("currentUser", session.CurrentUser(w, r))

	// If xml requested, serve with that template
	if strings.HasSuffix(r.URL.Path, ".xml") {
		view.Layout("")
		view.Template("stories/views/index.xml.got")
	}

	return view.Render()

}
