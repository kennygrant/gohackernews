package useractions

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/kennygrant/hackernews/src/lib/authorise"
	"github.com/kennygrant/hackernews/src/lib/status"
	"github.com/kennygrant/hackernews/src/users"
)

// HandleCreateShow handles GET users/create
func HandleCreateShow(context router.Context) error {

	// Authorise
	err := authorise.Path(context)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Setup
	view := view.New(context)
	user := users.New()
	view.AddKey("user", user)

	// Serve
	return view.Render()
}

// HandleCreate handles POST users/create
func HandleCreate(context router.Context) error {

	// Authorise
	err := authorise.Path(context)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Setup context
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err)
	}

	// We  check for duplicates in here - name and email must be unique
	count, err := users.Query().Where("email=?", params.Get("email")).Count()
	if err != nil {
		return router.InternalError(err)
	}
	if count > 0 {
		return router.NotAuthorizedError(err, "User already exists", "Sorry, a user already exists with that email.")
	}

	count, err = users.Query().Where("name=?", params.Get("name")).Count()
	if err != nil {
		return router.InternalError(err)
	}
	if count > 0 {
		return router.NotAuthorizedError(err, "User already exists", "Sorry, a user already exists with that name, please choose another.")
	}

	// Set some defaults for the new user
	params.SetInt("status", status.Published)
	params.SetInt("role", users.RoleReader)
	params.SetInt("points", 1)

	// Now try to create the user
	id, err := users.Create(params.Map())
	if err != nil {
		return router.InternalError(err, "Error", "Sorry, an error occurred creating the user record.")
	}

	context.Logf("#info Created user id,%d", id)

	// Find the user again so we can save login
	user, err := users.Find(id)
	if err != nil {
		context.Logf("#error parsing user id: %s", err)
		return router.NotFoundError(err)
	}

	// Save the fact user is logged in to session cookie
	err = loginUser(context, user)
	if err != nil {
		return router.InternalError(err)
	}

	// Redirect to root
	return router.Redirect(context, "/?message=welcome")
}
