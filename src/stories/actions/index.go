package storyactions

import (
	"fmt"
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

	windowTitle := context.Config("meta_title")
	switch filter {
	case "Video:":
		windowTitle = "Golang Videos"
	}

	// Render the template
	view := view.New(context)
	view.AddKey("page", page)
	view.AddKey("stories", results)
	view.AddKey("pubdate", storiesModTime(results))
	view.AddKey("meta_title", windowTitle)
	view.AddKey("meta_desc", context.Config("meta_desc"))
	view.AddKey("meta_keywords", context.Config("meta_keywords"))
	view.AddKey("meta_rss", storiesXMLPath(context))

	if context.Param("format") == ".xml" {
		view.Layout("")
		view.Template("stories/views/index.xml.got")
	}

	return view.Render()

}

// storiesModTime returns the mod time of the first story, or current time if no stories
func storiesModTime(availableStories []*stories.Story) time.Time {
	if len(availableStories) == 0 {
		return time.Now()
	}
	story := availableStories[0]

	return story.UpdatedAt
}

// storiesXMLPath returns the xml path for a given request to a stories link
func storiesXMLPath(context router.Context) string {

	request := context.Request()

	p := strings.Replace(request.URL.Path, ".xml", "", 1)
	if p == "/" {
		p = "/index"
	}

	q := request.URL.RawQuery
	if len(q) > 0 {
		q = "?" + q
	}

	return fmt.Sprintf("%s.xml%s", p, q)
}
