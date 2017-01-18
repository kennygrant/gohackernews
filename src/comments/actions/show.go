package commentactions

import (
	"net/http"

	"github.com/fragmenta/auth/can"
	"github.com/fragmenta/mux"
	"github.com/fragmenta/server"
	"github.com/fragmenta/view"

	"github.com/kennygrant/gohackernews/src/comments"
	"github.com/kennygrant/gohackernews/src/lib/session"
	"github.com/kennygrant/gohackernews/src/lib/status"
)

// HandleShow displays a single comment.
func HandleShow(w http.ResponseWriter, r *http.Request) error {

	// Fetch the  params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Find the comment
	comment, err := comments.Find(params.GetInt(comments.KeyName))
	if err != nil {
		return server.NotFoundError(err)
	}

	// Authorise access
	currentUser := session.CurrentUser(w, r)

	// Authorise access - for now all stories are visible, later might control on draft/published
	if comment.Status < status.None { // status.Published
		err = can.Show(comment, currentUser)
		if err != nil {
			return server.NotAuthorizedError(err)
		}
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.CacheKey(comment.CacheKey())
	view.AddKey("comment", comment)
	view.AddKey("currentUser", currentUser)
	return view.Render()
}
