package storyactions

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/kennygrant/hackernews/src/comments"
	"github.com/kennygrant/hackernews/src/lib/authorise"
	"github.com/kennygrant/hackernews/src/stories"
)

// HandleShow displays a single story
func HandleShow(context router.Context) error {

	// Find the story
	story, err := stories.Find(context.ParamInt("id"))
	if err != nil {
		return router.InternalError(err)
	}

	// Find the comments for this story
	// Fetch the comments
	q := comments.Where("story_id=?", story.Id).Order(comments.RankOrder)
	rootComments, err := comments.FindAll(q)
	if err != nil {
		return router.InternalError(err)
	}

	// Render the template
	view := view.New(context)
	view.AddKey("story", story)
	view.AddKey("meta_title", story.Name)
	view.AddKey("meta_desc", story.Summary)
	view.AddKey("meta_keywords", story.Name)
	view.AddKey("comments", rootComments)
	view.AddKey("authenticity_token", authorise.CreateAuthenticityToken(context))

	return view.Render()
}
