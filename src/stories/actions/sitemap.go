package storyactions

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/kennygrant/gohackernews/src/stories"
)

// HandleSiteMap renders a site map of top stories
func HandleSiteMap(context router.Context) error {

	// Build a query
	q := stories.Query().Limit(5000)

	// Select only above 0 points,  Order by points, then id
	q.Where("points > 0").Order("points desc, id desc")

	// Fetch the stories
	results, err := stories.FindAll(q)
	if err != nil {
		return router.InternalError(err)
	}

	// Render the template
	view := view.New(context)
	view.Layout("")
	view.Template("stories/views/sitemap.xml.got")
	view.AddKey("stories", results)
	view.AddKey("pubdate", storiesModTime(results))
	return view.Render()
}
