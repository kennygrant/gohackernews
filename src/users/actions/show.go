package useractions

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/kennygrant/gohackernews/src/comments"
	"github.com/kennygrant/gohackernews/src/stories"
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

	// Get the user comments
	q := comments.Where("user_id=?", user.Id).Limit(10).Order("created_at desc")
	userComments, err := comments.FindAll(q)
	if err != nil {
		return router.InternalError(err)
	}

	// Get the user stories
	q = stories.Where("user_id=?", user.Id).Limit(50).Order("created_at desc")
	userStories, err := stories.FindAll(q)
	if err != nil {
		return router.InternalError(err)
	}

	// Render the Template
	view := view.New(context)
	view.AddKey("user", user)
	view.AddKey("comments", userComments)
	view.AddKey("stories", userStories)
	view.AddKey("meta_title", user.Name)
	view.AddKey("meta_desc", user.Name)
	return view.Render()

}
