package storyactions

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/kennygrant/gohackernews/src/lib/authorise"
	"github.com/kennygrant/gohackernews/src/stories"
)

const (
	gravity = 1.8
)

// HandleUpdateShow renders the form to update a story
func HandleUpdateShow(context router.Context) error {

	// Find the story
	story, err := stories.Find(context.ParamInt("id"))
	if err != nil {
		return router.NotFoundError(err)
	}

	// Authorise update story
	err = authorise.Resource(context, story)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Render the template
	view := view.New(context)
	view.AddKey("story", story)

	return view.Render()
}

// HandleUpdate handles the POST of the form to update a story
func HandleUpdate(context router.Context) error {

	// Find the story
	story, err := stories.Find(context.ParamInt("id"))
	if err != nil {
		return router.NotFoundError(err)
	}

	// Authorise update story
	err = authorise.ResourceAndAuthenticity(context, story)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Update the story from params
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err)
	}

	// Clean params according to role
	accepted := stories.AllowedParams()
	if authorise.CurrentUser(context).Admin() {
		accepted = stories.AllowedParamsAdmin()
	}
	cleanedParams := params.Clean(accepted)

	err = story.Update(cleanedParams)
	if err != nil {
		return err // Create returns a router.Error
	}

	err = updateStoriesRank()
	if err != nil {
		return router.InternalError(err)
	}

	// Redirect to story
	return router.Redirect(context, story.URLShow())
}
