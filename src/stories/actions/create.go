package storyactions

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/kennygrant/hackernews/src/lib/authorise"
	"github.com/kennygrant/hackernews/src/stories"
)

// HandleCreateShow serves the create form via GET for stories
func HandleCreateShow(context router.Context) error {

	// Authorise
	err := authorise.Path(context)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Render the template
	view := view.New(context)
	story := stories.New()
	view.AddKey("story", story)
	view.AddKey("meta_title", "Go Hacker News Submit")
	view.AddKey("authenticity_token", authorise.CreateAuthenticityToken(context))
	return view.Render()
}

// HandleCreate handles the POST of the create form for stories
func HandleCreate(context router.Context) error {

	// Check csrf token
	err := authorise.AuthenticityToken(context)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Check permissions - if not logged in and above 1 points, redirect to error
	if !authorise.CurrentUser(context).CanSubmit() {
		return router.NotAuthorizedError(nil, "Sorry", "You need to be registered and have more than 1 points to submit stories.")
	}

	// Get user details
	user := authorise.CurrentUser(context)
	ip := getUserIP(context)

	// Check that no story with this url already exists
	q := stories.Where("url=?", context.Param("url"))
	duplicates, err := stories.FindAll(q)
	if err != nil {
		return router.InternalError(err)
	}

	if len(duplicates) > 0 {
		dupe := duplicates[0]
		// Add a point to dupe
		addStoryVote(dupe, user, ip, 1)

		return router.Redirect(context, dupe.URLShow())
	}

	// Setup context
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err)
	}

	// Set a few params
	params.SetInt("points", 1)
	params.SetInt("user_id", user.Id)
	params.Set("user_name", user.Name)

	id, err := stories.Create(params.Map())
	if err != nil {
		return err // Create returns a router.Error
	}

	// Log creation
	context.Logf("#info Created story id,%d", id)

	// Redirect to the new story
	story, err := stories.Find(id)
	if err != nil {
		return router.InternalError(err)
	}

	// We need to add a vote to the story here too by adding a join to the new id
	err = recordStoryVote(story, user, ip, +1)
	if err != nil {
		return err
	}

	// Re-rank stories
	err = updateStoriesRank()
	if err != nil {
		return err
	}

	return router.Redirect(context, story.URLIndex())
}
