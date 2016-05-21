package storyactions

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/kennygrant/gohackernews/src/stories"
)

// HandleHome displays a list of stories using gravity to order them
// used for the home page for gravity rank see votes.go
// responds to GET /
func HandleHome(context router.Context) error {

	// Build a query
	q := stories.Query().Limit(listLimit)

	// Select only above 0 points,  Order by rank, then points, then name
	q.Where("points > 0").Order("rank desc, points desc, id desc")

	// Set the offset in pages if we have one
	page := int(context.ParamInt("page"))
	if page > 0 {
		q.Offset(listLimit * page)
	}

	// Fetch the stories
	results, err := stories.FindAll(q)
	if err != nil {
		return router.InternalError(err)
	}

	// Render the template
	view := view.New(context)
	setStoriesMetadata(view, context.Request())
	view.AddKey("page", page)
	view.AddKey("stories", results)
	view.Template("stories/views/index.html.got")

	if context.Param("format") == ".xml" {
		view.Layout("")
		view.Template("stories/views/index.xml.got")
	}

	return view.Render()

}
