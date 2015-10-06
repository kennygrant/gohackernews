package storyactions

import (
	"strings"

	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/kennygrant/gohackernews/src/lib/authorise"
	"github.com/kennygrant/gohackernews/src/stories"
)

const listLimit = 100

// HandleHome displays a list of stories using gravity to order them
// used for the home page for gravity rank see votes.go
func HandleHome(context router.Context) error {

	// Build a query
	q := stories.Query().Limit(listLimit)

	// Select only above 0 points,  Order by rank, then points, then name
	q.Where("points > 0").Order("rank desc, points desc, id desc")

	// Fetch the stories
	results, err := stories.FindAll(q)
	if err != nil {
		return router.InternalError(err)
	}

	// Render the template
	view := view.New(context)
	view.AddKey("stories", results)
	view.AddKey("meta_title", "Golang News")
	view.AddKey("meta_desc", "News for Go Hackers, in the style of Hacker News. A curated selection of the latest links about the Go programming language.")
	view.AddKey("meta_keywords", "golang news, blog, links, go developers, go web apps, web applications, fragmenta")
	view.AddKey("authenticity_token", authorise.CreateAuthenticityToken(context))

	view.Template("stories/views/index.html.got")
	return view.Render()

}

// HandleIndex displays a list of stories at /stories
func HandleIndex(context router.Context) error {

	// Build a query
	q := stories.Query().Limit(listLimit)

	// Order by date by default
	q.Where("points > -6").Order("created_at desc")

	// Filter if necessary - this assumes name and summary cols
	filter := context.Param("filter")
	if len(filter) > 0 {
		context.Logf("FILTER %s", filter)

		// Replace special characters with escaped sequence
		filter = strings.Replace(filter, "_", "\\_", -1)
		filter = strings.Replace(filter, "%", "\\%", -1)

		// initially very simple, do ilike query for filter with wildcards
		q.Where("stories.name ILIKE ?", "%"+filter+"%")

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
	view.AddKey("authenticity_token", authorise.CreateAuthenticityToken(context))

	return view.Render()

}
