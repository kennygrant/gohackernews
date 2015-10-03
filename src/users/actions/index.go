package useractions

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/kennygrant/gohackernews/src/lib/authorise"
	"github.com/kennygrant/gohackernews/src/users"
)

// HandleIndex serves a GET request at /users
func HandleIndex(context router.Context) error {

	// Authorise
	err := authorise.Path(context)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Fetch the users
	q := users.Query().Order("role desc, created_at desc")
	userList, err := users.FindAll(q)
	if err != nil {
		context.Logf("#error Error indexing users %s", err)
		return router.InternalError(err)
	}

	// Serve template
	view := view.New(context)
	view.AddKey("users", userList)
	return view.Render()

}
