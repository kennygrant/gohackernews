package commentactions

import (
	"fmt"

	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/kennygrant/gohackernews/src/comments"
	"github.com/kennygrant/gohackernews/src/lib/authorise"
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

	// Clean params according to role
	accepted := comments.AllowedParams()
	if authorise.CurrentUser(context).Admin() {
		accepted = comments.AllowedParamsAdmin()
	}
	cleanedParams := params.Clean(accepted)

	err = comment.Update(cleanedParams)
	if err != nil {
		return router.InternalError(err)
	}

	// Redirect to the story concerned, or the comment if no story
	if comment.StoryId > 0 {
		return router.Redirect(context, fmt.Sprintf("/stories/%d", comment.StoryId))
	}

	return router.Redirect(context, fmt.Sprintf("/comments/%d", comment.Id))
}
