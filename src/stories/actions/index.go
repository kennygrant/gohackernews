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
// responds to GET /
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
	setStoriesMetadata(view)
	view.AddKey("stories", results)
	view.AddKey("authenticity_token", authorise.CreateAuthenticityToken(context))
	view.Template("stories/views/index.html.got")
	return view.Render()

}

// HandleCode displays a list of stories linking to repos (github etc) using gravity to order them
// responds to GET /stories/code
func HandleCode(context router.Context) error {

	// Build a query
	q := stories.Query().Where("points > -6").Order("rank desc, points desc, id desc").Limit(listLimit)

	// Restrict to stories with have a url starting with github.com or bitbucket.org
	// other code repos can be added later
	q.Where("url ILIKE 'https://github.com%'").OrWhere("url ILIKE 'https://bitbucket.org'")

	// Fetch the stories
	results, err := stories.FindAll(q)
	if err != nil {
		return router.InternalError(err)
	}

	// Render the template
	view := view.New(context)
	setStoriesMetadata(view)
	view.AddKey("stories", results)
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
	filter := context.Param("q")
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

	// Fetch the stories
	results, err := stories.FindAll(q)
	if err != nil {
		return router.InternalError(err)
	}

	// Render the template
	view := view.New(context)
	setStoriesMetadata(view)
	view.AddKey("filter", filter)
	view.AddKey("stories", results)
	view.AddKey("meta_title", "Golang News links")
	view.AddKey("authenticity_token", authorise.CreateAuthenticityToken(context))

	return view.Render()

}

func setStoriesMetadata(view *view.Renderer) {
	view.AddKey("meta_title", "Golang News")
	view.AddKey("meta_desc", "News for Go Hackers, in the style of Hacker News. A curated selection of the latest links about the Go programming language.")
	view.AddKey("meta_keywords", "golang news, blog, links, go developers, go web apps, web applications, fragmenta")

}
