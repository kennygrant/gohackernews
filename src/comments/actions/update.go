package commentactions

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/kennygrant/hackernews/src/comments"
	"github.com/kennygrant/hackernews/src/lib/authorise"
)

// HandleUpdateShow responds to GET /comments/update with the form to update a comment
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
	view.AddKey("authenticity_token", authorise.CreateAuthenticityToken(context))
	return view.Render()
}

// HandleUpdate responds to POST /comments/update
func HandleUpdate(context router.Context) error {

	// Find the comment
	comment, err := comments.Find(context.ParamInt("id"))
	if err != nil {
		return router.NotFoundError(err)
	}

	// Authorise update comment, check auth token
	err = authorise.ResourceAndAuthenticity(context, comment)
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
