package commentactions

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/kennygrant/hackernews/src/comments"
	"github.com/kennygrant/hackernews/src/lib/authorise"
)

// HandleShow displays a single comment
func HandleShow(context router.Context) error {

	// Find the comment
	comment, err := comments.Find(context.ParamInt("id"))
	if err != nil {
		return router.InternalError(err)
	}

	// Authorise access
	err = authorise.Resource(context, comment)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Render the template
	view := view.New(context)
	view.AddKey("comment", comment)
	return view.Render()
}
