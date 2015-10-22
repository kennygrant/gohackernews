package useractions

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/kennygrant/gohackernews/src/lib/authorise"
	"github.com/kennygrant/gohackernews/src/users"
)

// HandleUpdateShow serves a get request at /users/1/update (show form to update)
func HandleUpdateShow(context router.Context) error {
	// Setup context for template
	view := view.New(context)

	user, err := users.Find(context.ParamInt("id"))
	if err != nil {
		context.Logf("#error Error finding user %s", err)
		return router.NotFoundError(err)
	}

	// Authorise
	err = authorise.Resource(context, user)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	view.AddKey("user", user)

	return view.Render()
}

// HandleUpdate or PUT /users/1/update
func HandleUpdate(context router.Context) error {

	// Find the user
	id := context.ParamInt("id")
	user, err := users.Find(id)
	if err != nil {
		context.Logf("#error Error finding user %s", err)
		return router.NotFoundError(err)
	}

	// Authorise
	err = authorise.ResourceAndAuthenticity(context, user)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Get the params
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err)
	}

	// Clean params according to role
	accepted := users.AllowedParams()
	if authorise.CurrentUser(context).Admin() {
		accepted = users.AllowedParamsAdmin()
	}
	allowedParams := params.Clean(accepted)
	err = user.Update(allowedParams)
	if err != nil {
		return router.InternalError(err)
	}

	// Redirect to user
	return router.Redirect(context, user.URLShow())
}
