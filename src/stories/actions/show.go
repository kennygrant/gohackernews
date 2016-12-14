package storyactions

import (
	"fmt"

	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/kennygrant/gohackernews/src/comments"
	"github.com/kennygrant/gohackernews/src/lib/stats"
	"github.com/kennygrant/gohackernews/src/stories"
)

// HandleShow displays a single story
func HandleShow(context router.Context) error {
	stats.RegisterHit(context)

	// Find the story
	story, err := stories.Find(context.ParamInt("id"))
	if err != nil {
		return router.InternalError(err)
	}

	// Redirect requests to the canonical url
	if context.Path() != story.CanonicalURL() {
		return router.Redirect(context, story.CanonicalURL())
	}

	// Find the comments for this story
	// Fetch the comments
	q := comments.Where("story_id=?", story.Id).Order(comments.RankOrder)
	rootComments, err := comments.FindAll(q)
	if err != nil {
		return router.InternalError(err)
	}

	meta := story.Summary
	if len(meta) == 0 {
		meta = fmt.Sprintf("A story on %s, %s", context.Config("meta_title"), context.Config("meta_desc"))
	}

	// Render the template
	view := view.New(context)
	view.AddKey("story", story)
	view.AddKey("meta_title", story.Name)
	view.AddKey("meta_desc", meta)
	view.AddKey("meta_keywords", fmt.Sprintf("%s %s", story.Name, context.Config("meta_keywords")))
	view.AddKey("comments", rootComments)

	return view.Render()
}
