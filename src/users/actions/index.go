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

	// Query for most recent 100 users
	q := users.Query().Order("created_at desc").Limit(100)

	// Fetch 100 of them
	userList, err := users.FindAll(q)
	if err != nil {
		return router.InternalError(err)
	}

	// Get a count of all users
	count, err := q.Count()
	if err != nil {
		return router.InternalError(err)
	}

	// Get a count of admin users
	adminsCount, err := q.Where("role=100").Count()
	if err != nil {
		return router.InternalError(err)
	}

	// Serve template
	view := view.New(context)
	view.AddKey("users", userList)
	view.AddKey("count", count)
	view.AddKey("adminsCount", adminsCount)
	return view.Render()

}
