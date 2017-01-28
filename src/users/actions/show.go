package useractions

import (
	"net/http"

	"github.com/fragmenta/auth/can"
	"github.com/fragmenta/mux"
	"github.com/fragmenta/server"
	"github.com/fragmenta/view"

	"github.com/kennygrant/gohackernews/src/lib/session"
	"github.com/kennygrant/gohackernews/src/users"
)

// HandleShow displays a single user.
func HandleShow(w http.ResponseWriter, r *http.Request) error {

	// Fetch the  params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Find the user
	user, err := users.Find(params.GetInt(users.KeyName))
	if err != nil {
		return server.NotFoundError(err)
	}

	currentUser := session.CurrentUser(w, r)

	// Authorise access
	err = can.Show(user, currentUser)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.CacheKey(user.CacheKey())
	view.AddKey("user", user)
	view.AddKey("currentUser", currentUser)
	return view.Render()
}

// HandleShowName redirects a GET request of /u/username to the user show page
func HandleShowName(w http.ResponseWriter, r *http.Request) error {

	// Fetch the  params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Find the user by name
	q := users.Where("name=?", params.Get("name"))
	results, err := users.FindAll(q)
	if err != nil {
		return server.NotFoundError(err, "Error finding user")
	}

	// If valid query but no results
	if len(results) == 0 {
		return server.NotFoundError(err, "User not found")
	}

	// Redirect to user show page
	return server.Redirect(w, r, results[0].ShowURL())
}
