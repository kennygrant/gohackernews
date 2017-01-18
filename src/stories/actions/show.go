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

	fmt.Printf("HERE:%d:\n", story.Status)

	// Authorise access - for now all stories are visible, later might control on draft/published
	if story.Status < status.None { // status.Published
		err = can.Show(story, session.CurrentUser(w, r))
		if err != nil {
			return server.NotAuthorizedError(err)
		}
	}

	// Find the comments for this story
	q := comments.Where("story_id=?", story.ID).Order(comments.Order)
	comments, err := comments.FindAll(q)
	if err != nil {
		return server.InternalError(err)
	}

	meta := story.Summary
	if len(meta) == 0 {
		meta = fmt.Sprintf("A story on %s, %s", config.Get("meta_title"), config.Get("meta_desc"))
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.CacheKey(story.CacheKey())
	view.AddKey("story", story)
	view.AddKey("meta_title", story.Name)
	view.AddKey("meta_desc", meta)
	view.AddKey("meta_keywords", fmt.Sprintf("%s %s", story.Name, config.Get("meta_keywords")))
	view.AddKey("comments", comments)
	view.AddKey("currentUser", session.CurrentUser(w, r))
	return view.Render()
}
