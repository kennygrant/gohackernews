package storyactions

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/kennygrant/gohackernews/src/stories"
)

// HandleCode displays a list of stories linking to repos (github etc) using gravity to order them
// responds to GET /stories/code
func HandleCode(context router.Context) error {

	// Build a query
	q := stories.Query().Where("points > -6").Order("rank desc, points desc, id desc").Limit(listLimit)

	// Restrict to stories with have a url starting with github.com or bitbucket.org
	// other code repos can be added later
	q.Where("url ILIKE 'https://github.com%'").OrWhere("url ILIKE 'https://bitbucket.org'")

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
