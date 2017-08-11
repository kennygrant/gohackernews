package storyactions

import (
	"fmt"
	"net/http"

	"github.com/fragmenta/auth/can"
	"github.com/fragmenta/mux"
	"github.com/fragmenta/server"
	"github.com/fragmenta/server/config"
	"github.com/fragmenta/view"

	"github.com/kennygrant/gohackernews/src/comments"
	"github.com/kennygrant/gohackernews/src/lib/session"
	"github.com/kennygrant/gohackernews/src/lib/status"
	"github.com/kennygrant/gohackernews/src/stories"
)

// HandleShow displays a single story.
func HandleShow(w http.ResponseWriter, r *http.Request) error {

	// Fetch the  params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Find the story
	story, err := stories.Find(params.GetInt(stories.KeyName))
	if err != nil {
		return server.NotFoundError(err)
	}

	// Authorise access - for now all stories are visible, later might control on draft/published
	if story.Status < status.None { // status.Published
		err = can.Show(story, session.CurrentUser(w, r))
		if err != nil {
			return server.NotAuthorizedError(err)
		}
	}

	// Find the comments for this story, excluding those under 0
	q := comments.Where("story_id=?", story.ID).Where("points > 0").Order(comments.Order)
	comments, err := comments.FindAll(q)
	if err != nil {
		return server.InternalError(err)
	}

	meta := story.Summary
	if meta == "" {
		meta = fmt.Sprintf("%s - %s", config.Get("meta_title"), config.Get("meta_desc"))
	} else if len(meta) < 50 {
		meta = fmt.Sprintf("%s - %s", meta, config.Get("meta_desc"))
	}
	metaTitle := fmt.Sprintf("%s - %s", story.Name, config.Get("meta_title"))

	// Render the template
	view := view.NewRenderer(w, r)
	view.CacheKey(story.CacheKey())
	view.AddKey("story", story)
	view.AddKey("meta_title", metaTitle)
	view.AddKey("meta_desc", meta)
	view.AddKey("meta_foot", config.Get("meta_desc"))
	view.AddKey("meta_keywords", fmt.Sprintf("%s %s", story.Name, config.Get("meta_keywords")))
	view.AddKey("comments", comments)
	view.AddKey("currentUser", session.CurrentUser(w, r))
	return view.Render()
}
