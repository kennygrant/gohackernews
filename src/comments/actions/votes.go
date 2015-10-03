package commentactions

import (
	"fmt"

	"github.com/fragmenta/query"
	"github.com/fragmenta/router"

	"github.com/kennygrant/gohackernews/src/comments"
	"github.com/kennygrant/gohackernews/src/lib/authorise"
	"github.com/kennygrant/gohackernews/src/users"
)

// HandleFlag handles POST to /comments/123/flag
func HandleFlag(context router.Context) error {
	// Find the comment
	comment, err := comments.Find(context.ParamInt("id"))
	if err != nil {
		return router.NotFoundError(err)
	}
	user := authorise.CurrentUser(context)
	ip := getUserIP(context)

	// Check we have no votes already from this user, if we do fail
	if commentHasUserFlag(comment, user) {
		return router.NotAuthorizedError(err, "Flag Failed", "Sorry you are not allowed to flag twice, nice try!")
	}

	// Authorise upvote on comment for this user - our rules are:
	if !user.CanFlag() {
		return router.NotAuthorizedError(err, "Flag Failed", "Sorry, you can't flag yet")
	}

	err = adjustUserPoints(user, -1)
	if err != nil {
		return err
	}

	err = addCommentVote(comment, user, ip, -5)
	if err != nil {
		return err
	}
	return updateCommentsRank(comment.StoryId)
}

// HandleDownvote handles POST to /comments/123/downvote
func HandleDownvote(context router.Context) error {
	// Find the comment
	comment, err := comments.Find(context.ParamInt("id"))
	if err != nil {
		return router.NotFoundError(err)
	}
	user := authorise.CurrentUser(context)
	ip := getUserIP(context)

	if !user.Admin() {
		// Check we have no votes already from this user, if we do fail
		if commentHasUserVote(comment, user) {
			return router.NotAuthorizedError(err, "Vote Failed", "Sorry you are not allowed to vote twice, nice try!")
		}
	}

	// Authorise upvote on comment for this user - our rules are:
	if !user.CanDownvote() {
		return router.NotAuthorizedError(err, "Vote Failed", "Sorry, you can't downvote yet")
	}

	err = adjustUserPoints(user, -1)
	if err != nil {
		return err
	}

	// Adjust points on comment and add to the vote table
	err = addCommentVote(comment, user, ip, -1)
	if err != nil {
		return err
	}

	return updateCommentsRank(comment.StoryId)
}

// HandleUpvote handles POST to /comments/123/upvote
func HandleUpvote(context router.Context) error {

	// Find the comment
	comment, err := comments.Find(context.ParamInt("id"))
	if err != nil {
		return router.NotFoundError(err)
	}

	user := authorise.CurrentUser(context)
	ip := getUserIP(context)

	if !user.Admin() {
		// Check we have no votes already from this user, if we do fail
		if commentHasUserVote(comment, user) {
			return router.NotAuthorizedError(err, "Vote Failed", "Sorry you are not allowed to vote twice, nice try!")
		}
	}

	// Authorise upvote on comment for this user - our rules are:
	if !user.CanUpvote() {
		return router.NotAuthorizedError(err, "Vote Failed", "Sorry, you can't upvote yet")
	}

	// Adjust points on comment and add to the vote table
	err = addCommentVote(comment, user, ip, +1)
	if err != nil {
		return err
	}

	return updateCommentsRank(comment.StoryId)
}

// addCommentVote adjusts the comment points, and adds a vote record for this user
func addCommentVote(comment *comments.Comment, user *users.User, ip string, delta int64) error {

	if comment.Points < -5 && delta < 0 {
		return router.InternalError(nil, "Vote Failed", "Comment is already hidden")
	}

	// Update the comment points by delta
	err := comment.Update(map[string]string{"points": fmt.Sprintf("%d", comment.Points+delta)})
	if err != nil {
		return router.InternalError(err, "Vote Failed", "Sorry your adjust vote points")
	}

	// Update the *comment* user points by delta
	commentUser, err := users.Find(comment.UserId)
	err = adjustUserPoints(commentUser, +1)
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
		return router.InternalError(err, "Vote Failed", "Sorry could not adjust user points")
	}

	return nil
}

// recordCommentVote adds a vote record for this user
func recordCommentVote(comment *comments.Comment, user *users.User, ip string, delta int64) error {

	// Add an entry in the votes table
	// FIXME: adjust query to do this for us we should use ?,?,? here...
	// $1, $2 is surprising, shouldn't we expect query package to deal with this for us?
	_, err := query.Exec("insert into votes VALUES(now(),$1,NULL,$2,$3,$4)", comment.Id, user.Id, ip, delta)
	if err != nil {
		return router.InternalError(err, "Vote Failed", "Sorry your vote failed to record")
	}

	return nil
}

// commentHasUserVote returns true if we already have a vote for this comment from this user
func commentHasUserVote(comment *comments.Comment, user *users.User) bool {
	// Query votes table for rows with userId and commentId
	// if we don't get error, return true
	results, err := query.New("votes", "comment_id").Where("comment_id=?", comment.Id).Where("user_id=?", user.Id).Results()

	if err == nil && len(results) == 0 {
		return false
	}

	return true
}

// commentHasUserFlag returns true if we already have a flag for this comment from this user
func commentHasUserFlag(comment *comments.Comment, user *users.User) bool {
	// Query flags table for rows with userId and commentId
	// if we don't get error, return true
	results, err := query.New("flags", "comment_id").Where("comment_id=?", comment.Id).Where("user_id=?", user.Id).Results()
	if err == nil && len(results) == 0 {
		return false
	}

	return true
}

func updateCommentsRank(storyID int64) error {
	sql := "update comments set rank = points / POWER((select max(id) from comments) - id + 1,1.8) where story_id=$1"
	_, err := query.Exec(sql, storyID)
	return err
}

func getUserIP(router.Context) string {
	return ""
}
