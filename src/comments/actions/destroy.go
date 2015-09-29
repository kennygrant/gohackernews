package commentactions

import (
	"github.com/fragmenta/router"

	"github.com/kennygrant/hackernews/src/comments"
	"github.com/kennygrant/hackernews/src/lib/authorise"
)

// HandleDestroy handles a DESTROY request for comments
func HandleDestroy(context router.Context) error {

	// Find the comment
	comment, err := comments.Find(context.ParamInt("id"))
	if err != nil {
		return router.NotFoundError(err)
	}

	// Authorise destroy comment
	err = authorise.Resource(context, comment)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Destroy the comment
	comment.Destroy()

	// Redirect to comments root
	return router.Redirect(context, comment.URLIndex())
}
