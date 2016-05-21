package storyactions

import (
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/kennygrant/gohackernews/src/comments"
	"github.com/kennygrant/gohackernews/src/stories"
)

// HandleShow displays a single story
func HandleShow(context router.Context) error {

	// Find the story
	story, err := stories.Find(context.ParamInt("id"))
	if err != nil {
		return router.InternalError(err)
	}

	// Redirect requests to the canonical url
	if context.Path() != story.URLShow() {
		return router.Redirect(context, story.URLShow())
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

	return view.Render()
}
