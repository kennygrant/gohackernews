package useractions

import (
	"strings"

	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/kennygrant/gohackernews/src/lib/authorise"
	"github.com/kennygrant/gohackernews/src/lib/status"
	"github.com/kennygrant/gohackernews/src/users"
)

// HandleCreateShow handles GET /users/create
func HandleCreateShow(context router.Context) error {

	// No auth as anyone can create users in this app

	// Setup
	view := view.New(context)
	user := users.New()
	view.AddKey("user", user)

	// Serve
	return view.Render()
}

// HandleCreate handles POST /users/create from the register page
func HandleCreate(context router.Context) error {

	// Check csrf token
	err := authorise.AuthenticityToken(context)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Setup context
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err)
	}

	// Check for email duplicates
	email := params.Get("email")
	if len(email) > 0 {

		if len(email) < 3 || !strings.Contains(email, "@") {
			return router.InternalError(err, "Invalid email", "Please just miss out the email field, or use a valid email.")
		}

		count, err := users.Query().Where("email=?", email).Count()
		if err != nil {
			return router.InternalError(err)
		}
		if count > 0 {
			return router.NotAuthorizedError(err, "User already exists", "Sorry, a user already exists with that email.")
		}
	}

	// Check for invalid or duplicate names
	name := params.Get("name")
	if len(name) < 2 {
		return router.InternalError(err, "Name too short", "Please choose a username longer than 2 characters")
	}

	count, err := users.Query().Where("name=?", name).Count()
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
	id, err := users.Create(params.Clean(users.AllowedParams()))
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
