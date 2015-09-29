package users

import ()

// User privs at karma points
// for now these are relaxed
/*
  1 points - submit, comment
  10 points - upvote (they start with 10 points)
  50 points - downvote
  100 points - flag

	karma is collected for comment upvotes *only* not for story upvotes

  flagging and downvoting sacrifices one point from user
*/

// CanUpvote returns true if this user can upvote - just let all users upvote for now
func (u *User) CanUpvote() bool {
	return u.Points > 0
}

// CanDownvote returns true if this user can downvote
func (u *User) CanDownvote() bool {
	return u.Points > 20
}

// CanFlag returns true if this user can flag
func (u *User) CanFlag() bool {
	return u.Points > 10
}

// CanSubmit returns true if this user can submit
func (u *User) CanSubmit() bool {
	return u.Points > 0
}

// CanComment returns true if this user can comment
func (u *User) CanComment() bool {
	return u.Points > 0
}
