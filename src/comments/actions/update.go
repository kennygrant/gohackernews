package commentactions

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/kennygrant/hackernews/src/comments"
	"github.com/kennygrant/hackernews/src/lib/authorise"
)

// HandleUpdateShow renders the form to update a comment
func HandleUpdateShow(context router.Context) error {

	// Find the comment
	comment, err := comments.Find(context.ParamInt("id"))
	if err != nil {
		return router.NotFoundError(err)
	}

	// Authorise update comment
	err = authorise.Resource(context, comment)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Render the template
	view := view.New(context)
	view.AddKey("comment", comment)
	// view.AddKey("csrf",auth.CSRFToken(""))
	return view.Render()
}

// HandleUpdateShow handles the POST of the form to update a comment
func HandleUpdate(context router.Context) error {

	// Find the comment
	comment, err := comments.Find(context.ParamInt("id"))
	if err != nil {
		return router.NotFoundError(err)
	}

	// Authorise update comment
	err = authorise.Resource(context, comment)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Update the comment from params
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err)
	}
	err = comment.Update(params.Map())
	if err != nil {
		return router.InternalError(err)
	}

	// Redirect to comment
	return router.Redirect(context, comment.URLShow())
}
