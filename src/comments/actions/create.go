package commentactions

import (
	"fmt"

	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/kennygrant/gohackernews/src/comments"
	"github.com/kennygrant/gohackernews/src/lib/authorise"
	"github.com/kennygrant/gohackernews/src/stories"
)

// HandleCreateShow serves the create form via GET for comments
func HandleCreateShow(context router.Context) error {

	// Authorise
	err := authorise.Path(context)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Render the template
	view := view.New(context)
	comment := comments.New()
	view.AddKey("comment", comment)

	// TODO: May have to validate parent_id or story_id
	view.AddKey("story_id", context.ParamInt("story_id"))
	view.AddKey("parent_id", context.ParamInt("parent_id"))
	view.AddKey("authenticity_token", authorise.CreateAuthenticityToken(context))
	return view.Render()
}

// HandleCreate handles the POST of the create form for comments
func HandleCreate(context router.Context) error {

	// Authorise csrf token
	err := authorise.AuthenticityToken(context)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Check permissions - if not logged in and above 0 points, redirect
	if !authorise.CurrentUser(context).CanComment() {
		return router.NotAuthorizedError(nil, "Sorry", "You need to be registered and have more than 0 points to comment.")
	}

	// Setup context
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err)
	}

	// Find parent story - this must exist
	story, err := stories.Find(params.GetInt("story_id"))
	if err != nil {
		return router.NotFoundError(err)
	}

	params.SetInt("story_id", story.Id)
	params.Set("story_name", story.Name)

	// Set a few params
	user := authorise.CurrentUser(context)
	params.SetInt("user_id", user.Id)
	params.Set("user_name", user.Name)
	params.SetInt("points", 1)

	// Find the parent and set dotted id
	// these are of the form xx.xx. with a trailing dot
	// this saves us from saving twice on create
	parentID := context.ParamInt("parent_id")
	if parentID > 0 {
		parent, err := comments.Find(parentID)
		if err != nil {
			return router.NotFoundError(err)
		}
		context.Logf("PARENT:%d - %s", parent.Id, parent.DottedIds)
		params.Set("dotted_ids", fmt.Sprintf(parent.DottedIds+"."))
	}

	id, err := comments.Create(params.Map())
	if err != nil {
		return router.InternalError(err)
	}

	// Log creation
	context.Logf("#info Created comment id,%d", id)

	// Update the story comment Count

	storyParams := map[string]string{"comment_count": fmt.Sprintf("%d", story.CommentCount+1)}
	err = story.Update(storyParams)
	if err != nil {
		return router.InternalError(err, "Error", "Could not update story.")
	}

	// Redirect to the new comment
	m, err := comments.Find(id)
	if err != nil {
		return router.InternalError(err)
	}

	// Re-rank comments on this story
	err = updateCommentsRank(m.StoryId)
	if err != nil {
		return err
	}

	return router.Redirect(context, m.URLStory())
}
