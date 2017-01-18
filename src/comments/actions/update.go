package commentactions

import (
	"net/http"

	"github.com/fragmenta/auth/can"
	"github.com/fragmenta/mux"
	"github.com/fragmenta/server"
	"github.com/fragmenta/view"

	"github.com/kennygrant/gohackernews/src/comments"
	"github.com/kennygrant/gohackernews/src/lib/session"
)

// HandleUpdateShow renders the form to update a comment.
func HandleUpdateShow(w http.ResponseWriter, r *http.Request) error {

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

	// Authorise update comment
	err = can.Update(comment, session.CurrentUser(w, r))
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.AddKey("comment", comment)
	return view.Render()
}

// HandleUpdate handles the POST of the form to update a comment
func HandleUpdate(w http.ResponseWriter, r *http.Request) error {

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

	// Check the authenticity token
	err = session.CheckAuthenticity(w, r)
	if err != nil {
		return err
	}

	// Authorise update comment
	err = can.Update(comment, session.CurrentUser(w, r))
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Validate the params, removing any we don't accept
	commentParams := comment.ValidateParams(params.Map(), comments.AllowedParams())

	err = comment.Update(commentParams)
	if err != nil {
		return server.InternalError(err)
	}

	// Redirect to comment
	return server.Redirect(w, r, comment.ShowURL())
}
