package useractions

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/kennygrant/gohackernews/src/users"
)

// HandleShow serve a get request at /users/1
func HandleShow(context router.Context) error {

	// No auth - this is public

	// Find the user
	user, err := users.Find(context.ParamInt("id"))
	if err != nil {
		context.Logf("#error parsing user id: %s", err)
		return router.NotFoundError(err)
	}

	// Set up view
	view := view.New(context)

	// Render the Template
	view.AddKey("user", user)
	view.AddKey("meta_title", user.Name)
	view.AddKey("meta_desc", user.Name)
	return view.Render()

}
