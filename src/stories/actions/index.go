package storyactions

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/kennygrant/gohackernews/src/stories"
)

const listLimit = 100

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

	if context.Param("format") == ".xml" {
		view.Layout("")
		view.Template("stories/views/index.xml.got")
	}

	return view.Render()

}

func setStoriesMetadata(view *view.Renderer, request *http.Request) {
	view.AddKey("pubdate", time.Now()) // could use latest story date instead?
	view.AddKey("meta_title", "Golang News")
	view.AddKey("meta_desc", "News for Go Hackers, in the style of Hacker News. A curated selection of the latest links about the Go programming language.")
	view.AddKey("meta_keywords", "golang news, blog, links, go developers, go web apps, web applications, fragmenta")

	p := strings.Replace(request.URL.Path, ".xml", "", 1)
	if p == "/" {
		p = "/index"
	}

	q := request.URL.RawQuery
	if len(q) > 0 {
		q = "?" + q
	}

	url := fmt.Sprintf("%s.xml%s", p, q)
	view.AddKey("meta_rss", url)

}
