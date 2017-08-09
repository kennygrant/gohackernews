package commentactions

import (
	"fmt"
	"net/http"

	"github.com/fragmenta/auth/can"
	"github.com/fragmenta/mux"
	"github.com/fragmenta/server"
	"github.com/fragmenta/server/log"
	"github.com/fragmenta/view"

	"github.com/kennygrant/gohackernews/src/comments"
	"github.com/kennygrant/gohackernews/src/lib/session"
	"github.com/kennygrant/gohackernews/src/stories"
)

// HandleCreateShow serves the create form via GET for comments.
func HandleCreateShow(w http.ResponseWriter, r *http.Request) error {

	comment := comments.New()

	// Authorise
	err := can.Create(comment, session.CurrentUser(w, r))
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.AddKey("comment", comment)
	return view.Render()
}

// HandleCreate handles the POST of the create form for comments
func HandleCreate(w http.ResponseWriter, r *http.Request) error {

	comment := comments.New()

	// Check the authenticity token
	err := session.CheckAuthenticity(w, r)
	if err != nil {
		return err
	}

	// Authorise
	currentUser := session.CurrentUser(w, r)
	err = can.Create(comment, currentUser)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Check permissions - if not logged in and above 0 points, redirect
	if !currentUser.CanComment() {
		return server.NotAuthorizedError(nil, "Sorry", "You need to be registered and have more than 0 points to comment.")
	}

	// Get Params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	text := params.Get("text")

	// Disallow empty comments
	if len(text) < 5 {
		return server.NotAuthorizedError(nil, "Comment too short", "Your comment is too short. Please try to post substantive comments which others will find useful.")
	}

	// Disallow comments which are too long
	if len(text) > 5000 {
		return server.NotAuthorizedError(nil, "Comment too long", "Your comment is too long.")
	}

	// Find parent story - this must exist
	story, err := stories.Find(params.GetInt("story_id"))
	if err != nil {
		return server.NotFoundError(err)
	}

	// Clean params according to role
	accepted := comments.AllowedParams()
	if currentUser.Admin() {
		accepted = comments.AllowedParamsAdmin()
	}
	commentParams := comment.ValidateParams(params.Map(), accepted)

	// Find the parent to set dotted id
	// these are of the form xx.xx. with a trailing dot
	// this saves us from saving twice on create
	parentID := params.GetInt("parent_id")
	if parentID > 0 {

		parent, e := comments.Find(parentID)
		if e != nil {
			return server.NotFoundError(err)
		}
		commentParams["dotted_ids"] = fmt.Sprintf(parent.DottedIDs + ".")
	}

	// Set other params from story/user details
	commentParams["story_id"] = fmt.Sprintf("%d", story.ID)
	commentParams["story_name"] = story.Name
	commentParams["user_id"] = fmt.Sprintf("%d", currentUser.ID)
	commentParams["user_name"] = currentUser.Name
	commentParams["points"] = "1"

	ID, err := comment.Create(commentParams)
	if err != nil {
		return server.InternalError(err)
	}

	// Log comment creation
	log.Info(log.Values{"msg": "Created comment", "comment_id": ID, "params": commentParams})

	// Update the story comment count with a count of comments above 0 points
	err = updateStoryCommentCount(story)
	if err != nil {
		return server.InternalError(err, "Error", "Could not update story.")
	}

	// Redirect to the new comment
	comment, err = comments.Find(ID)
	if err != nil {
		return server.InternalError(err)
	}

	// Re-rank comments on this story
	err = updateCommentsRank(comment.StoryID)
	if err != nil {
		return err
	}

	return server.Redirect(w, r, comment.StoryURL())
}
