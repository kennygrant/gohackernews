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
	currentUser := session.CurrentUser(w, r)
	err = can.Update(comment, currentUser)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.AddKey("currentUser", currentUser)
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
	currentUser := session.CurrentUser(w, r)
	err = can.Update(comment, currentUser)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Clean params according to role
	accepted := comments.AllowedParams()
	if currentUser.Admin() {
		accepted = comments.AllowedParamsAdmin()
	}
	commentParams := comment.ValidateParams(params.Map(), accepted)

	err = comment.Update(commentParams)
	if err != nil {
		return server.InternalError(err)
	}

	// Redirect to comment
	return server.Redirect(w, r, comment.ShowURL())
}
