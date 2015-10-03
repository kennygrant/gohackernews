package commentactions

import (
	"strings"

	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/kennygrant/hackernews/src/comments"
)

// HandleIndex displays a list of comments
func HandleIndex(context router.Context) error {

	// No auth this is public

	// Build a query to fetch latest 100 comments
	q := comments.Query().Limit(100).Order("created_at desc")

	// Filter on user id - we only show the actual user's comments
	// so not a nested view as in HN
	userID := context.ParamInt("u")
	if userID > 0 {
		q.Where("user_id=?", userID)
	}

	// Filter if necessary - this assumes name and summary cols
	filter := context.Param("filter")
	if len(filter) > 0 {
		filter = strings.Replace(filter, "&", "", -1)
		filter = strings.Replace(filter, " ", "", -1)
		filter = strings.Replace(filter, " ", " & ", -1)
		q.Where("(to_tsvector(text) @@ to_tsquery(?) )", filter)
	}

	// Fetch the comments
	results, err := comments.FindAll(q)
	if err != nil {
		return router.InternalError(err)
	}

	// Render the template
	view := view.New(context)
	view.AddKey("filter", filter)
	view.AddKey("comments", results)
	view.AddKey("meta_title", "Comments")
	return view.Render()

}
