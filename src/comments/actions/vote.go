package commentactions

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/fragmenta/mux"
	"github.com/fragmenta/query"
	"github.com/fragmenta/server"

	"github.com/kennygrant/gohackernews/src/comments"
	"github.com/kennygrant/gohackernews/src/lib/session"
	"github.com/kennygrant/gohackernews/src/stories"
	"github.com/kennygrant/gohackernews/src/users"
)

// HandleFlag handles POST to /comments/123/flag
func HandleFlag(w http.ResponseWriter, r *http.Request) error {

	// Check the authenticity token
	err := session.CheckAuthenticity(w, r)
	if err != nil {
		return err
	}

	// Fetch the  params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Find the comment
	comment, err := comments.Find(params.GetInt("id"))
	if err != nil {
		return server.NotFoundError(err)
	}
	user := session.CurrentUser(w, r)
	ip := getUserIP(r)

	// Check we have no votes already from this user, if we do fail
	if commentHasUserFlag(comment, user) {
		return server.NotAuthorizedError(err, "Flag Failed", "Sorry you are not allowed to flag twice, nice try!")
	}

	// Authorise upvote on comment for this user - our rules are:
	if !user.CanFlag() {
		return server.NotAuthorizedError(err, "Flag Failed", "Sorry, you can't flag yet")
	}

	// CURRENT User burns points for flagging
	err = adjustUserPoints(user, -1)
	if err != nil {
		return err
	}

	// Adjust the comment vote
	err = addCommentVote(comment, user, ip, -5)
	if err != nil {
		return err
	}

	err = updateCommentsRank(comment.StoryID)
	if err != nil {
		return err
	}

	// Adjust the story comment count
	story, err := stories.Find(comment.StoryID)
	if err != nil {
		return err
	}
	err = updateStoryCommentCount(story)
	if err != nil {
		return err
	}

	// Redirect to story
	return server.Redirect(w, r, fmt.Sprintf("/stories/%d", comment.StoryID))
}

// HandleDownvote handles POST to /comments/123/downvote
func HandleDownvote(w http.ResponseWriter, r *http.Request) error {

	// Check the authenticity token
	err := session.CheckAuthenticity(w, r)
	if err != nil {
		return err
	}

	// Fetch the  params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Find the comment
	comment, err := comments.Find(params.GetInt("id"))
	if err != nil {
		return server.NotFoundError(err)
	}
	user := session.CurrentUser(w, r)
	ip := getUserIP(r)

	if !user.Admin() {
		// Check we have no votes already from this user, if we do fail
		if commentHasUserVote(comment, user) {
			return server.NotAuthorizedError(err, "Vote Failed", "Sorry you are not allowed to vote twice, nice try!")
		}
	}

	// Authorise upvote on comment for this user - our rules are:
	if !user.CanDownvote() {
		return server.NotAuthorizedError(err, "Vote Failed", "Sorry, you can't downvote yet")
	}

	// CURRENT User burns points for downvoting
	err = adjustUserPoints(user, -1)
	if err != nil {
		return err
	}

	// Adjust points on comment and add to the vote table
	err = addCommentVote(comment, user, ip, -1)
	if err != nil {
		return err
	}

	// Adjust the story comment count if points now less than 1
	if comment.Points <= 1 {
		story, err := stories.Find(comment.StoryID)
		if err != nil {
			return err
		}
		err = updateStoryCommentCount(story)
		if err != nil {
			return err
		}
	}

	return updateCommentsRank(comment.StoryID)
}

// HandleUpvote handles POST to /comments/123/upvote
func HandleUpvote(w http.ResponseWriter, r *http.Request) error {

	// Check the authenticity token
	err := session.CheckAuthenticity(w, r)
	if err != nil {
		return err
	}

	// Fetch the  params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Find the comment
	comment, err := comments.Find(params.GetInt("id"))
	if err != nil {
		return server.NotFoundError(err)
	}

	user := session.CurrentUser(w, r)
	ip := getUserIP(r)

	if !user.Admin() {
		// Check we have no votes already from this user, if we do fail
		if commentHasUserVote(comment, user) {
			return server.NotAuthorizedError(err, "Vote Failed", "Sorry you are not allowed to vote twice, nice try!")
		}
	}

	// Authorise upvote on comment for this user - our rules are:
	if !user.CanUpvote() {
		return server.NotAuthorizedError(err, "Vote Failed", "Sorry, you can't upvote yet")
	}

	// Adjust points on comment and add to the vote table
	err = addCommentVote(comment, user, ip, +1)
	if err != nil {
		return err
	}

	return updateCommentsRank(comment.StoryID)
}

// addCommentVote adjusts the comment points, and adds a vote record for this user
func addCommentVote(comment *comments.Comment, user *users.User, ip string, delta int64) error {

	if comment.Points < -5 && delta < 0 {
		return server.InternalError(nil, "Vote Failed", "Comment is already hidden")
	}

	// Update the comment points by delta
	err := comment.Update(map[string]string{"points": fmt.Sprintf("%d", comment.Points+delta)})
	if err != nil {
		return server.InternalError(err, "Vote Failed", "Sorry your adjust vote points")
	}

	// Update the *comment* user points by delta
	commentUser, err := users.Find(comment.UserID)
	if err != nil {
		return err
	}
	err = adjustUserPoints(commentUser, delta)
	if err != nil {
		return err
	}

	return recordCommentVote(comment, user, ip, delta)
}

// removeUserPoints removes these points from this user
func adjustUserPoints(user *users.User, delta int64) error {

	// Update the user points
	err := user.Update(map[string]string{"points": fmt.Sprintf("%d", user.Points+delta)})
	if err != nil {
		return server.InternalError(err, "Vote Failed", "Sorry could not adjust user points")
	}

	return nil
}

// recordCommentVote adds a vote record for this user
func recordCommentVote(comment *comments.Comment, user *users.User, ip string, delta int64) error {

	// Add an entry in the votes table
	// FIXME: adjust query to do this for us we should use ?,?,? here...
	// $1, $2 is surprising, shouldn't we expect query package to deal with this for us?
	_, err := query.Exec("insert into votes VALUES(now(),$1,NULL,$2,$3,$4)", comment.ID, user.ID, ip, delta)
	if err != nil {
		return server.InternalError(err, "Vote Failed", "Sorry your vote failed to record")
	}

	return nil
}

// commentHasUserVote returns true if we already have a vote for this comment from this user
func commentHasUserVote(comment *comments.Comment, user *users.User) bool {
	// Query votes table for rows with userId and commentId
	// if we don't get error, return true
	results, err := query.New("votes", "comment_id").Where("comment_id=?", comment.ID).Where("user_id=?", user.ID).Results()

	if err == nil && len(results) == 0 {
		return false
	}

	return true
}

// commentHasUserFlag returns true if we already have a flag for this comment from this user
func commentHasUserFlag(comment *comments.Comment, user *users.User) bool {
	// Query flags table for rows with userId and commentId
	// if we don't get error, return true
	results, err := query.New("flags", "comment_id").Where("comment_id=?", comment.ID).Where("user_id=?", user.ID).Results()
	if err == nil && len(results) == 0 {
		return false
	}

	return true
}

// updateCommentsRank updates the rank of comments on this story
func updateCommentsRank(storyID int64) error {
	sql := "update comments set rank = 100 * points / POWER((select max(id) from comments) - id + 1,1.2) where story_id=$1"
	_, err := query.Exec(sql, storyID)
	return err
}

// updateStoryCommentCount updates a story for new comment counts
// discounting comments under 0 points
func updateStoryCommentCount(story *stories.Story) error {
	commentCount, err := comments.Query().Where("story_id=?", story.ID).Where("points > 0").Count()
	if err != nil {
		return err
	}
	storyParams := map[string]string{"comment_count": fmt.Sprintf("%d", commentCount)}
	return story.Update(storyParams)
}

func getUserIP(r *http.Request) string {
	// Store a hash of the ip (should we strip port?)
	ip := r.RemoteAddr
	forward := r.Header.Get("X-Forwarded-For")
	if len(forward) > 0 {
		ip = forward
	}

	// Hash for anonymity in our store
	hasher := sha256.New()
	hasher.Write([]byte(ip))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
