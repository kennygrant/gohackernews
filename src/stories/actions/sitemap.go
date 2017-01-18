package storyactions

import (
	"net/http"

	"github.com/fragmenta/server"
	"github.com/fragmenta/view"

	"github.com/kennygrant/gohackernews/src/stories"
)

// HandleSiteMap renders a site map of top stories
func HandleSiteMap(w http.ResponseWriter, r *http.Request) error {

	// Build a query
	q := stories.Query().Limit(5000)

	// Select only above 0 points,  Order by points, then id
	q.Where("points > 0").Order("points desc, id desc")

	// Fetch the stories
	results, err := stories.FindAll(q)
	if err != nil {
		return server.InternalError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.Layout("")
	view.Template("stories/views/sitemap.xml.got")
	view.AddKey("stories", results)
	view.AddKey("pubdate", storiesModTime(results))
	return view.Render()
}
