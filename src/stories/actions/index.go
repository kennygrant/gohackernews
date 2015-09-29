package storyactions

import (
	"strings"

	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/kennygrant/hackernews/src/lib/authorise"
	"github.com/kennygrant/hackernews/src/stories"
)

// HandleHome displays a list of stories using points and gravity to order them
// (with an optional filter perhaps?)
func HandleHome(context router.Context) error {

	// Authorise
	err := authorise.Path(context)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Build a query
	q := stories.Query().Limit(500)

	// Order by rank, then points, then name
	q.Order("rank desc, points desc, id desc")

	// Fetch the stories
	results, err := stories.FindAll(q)
	if err != nil {
		return router.InternalError(err)
	}

	// Render the template
	view := view.New(context)
	view.AddKey("stories", results)
	view.AddKey("meta_title", "Go Hacker News")
	view.AddKey("meta_desc", "News for golang Hackers, in the style of Hacker News")
	view.AddKey("meta_keywords", "golang news, blog, links, go developers")
	view.Template("stories/views/index.html.got")
	return view.Render()

}

// HandleIndex displays a list of stories
func HandleIndex(context router.Context) error {

	// Authorise
	err := authorise.Path(context)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Build a query
	q := stories.Query().Limit(500)

	// Order by date by default
	q.Order("created_at desc")

	// Filter if necessary - this assumes name and summary cols
	filter := context.Param("filter")
	if len(filter) > 0 {
		context.Logf("FILTER %s", filter)

		// Replace special characters with escaped sequence
		filter = strings.Replace(filter, "_", "\\_", -1)
		filter = strings.Replace(filter, "%", "\\%", -1)

		filter = "%" + filter + "%"

		// initially very simple, do ilike query
		q.Where("stories.name SIMILAR TO ?", filter)

		// If filtering, order by rank, not by date
		q.Order("rank desc, points desc, id desc")
	}

	// Fetch the stories
	results, err := stories.FindAll(q)
	if err != nil {
		return router.InternalError(err)
	}

	// Render the template
	view := view.New(context)
	view.AddKey("filter", filter)
	view.AddKey("stories", results)
	view.AddKey("meta_title", "Go Hacker News Links")
	return view.Render()

}
