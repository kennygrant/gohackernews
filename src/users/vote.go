package users

/*
  Proposed voting rules - for now these are relaxed
  1 points - submit, comment
  10 points - upvote (they start with 10 points)
  50 points - downvote
  100 points - flag

	karma is collected for comment upvotes *only* not for story upvotes
	karma is sacrificed in negative actions - flagging and downvoting
*/

// CanSubmit returns true if this user can submit.
func (u *User) CanSubmit() bool {
	return u.Points > 0
}

// CanComment returns true if this user can comment.
func (u *User) CanComment() bool {
	return u.Points > 0
}

// CanUpvote returns true if this user can upvote.
func (u *User) CanUpvote() bool {
	return u.Points > 2
}

// CanDownvote returns true if this user can downvote.
func (u *User) CanDownvote() bool {
	return u.Points > 20
}

// CanStyle returns true if this user can style text in comments/stories.
func (u *User) CanStyle() bool {
	return u.Points > 30
}

// CanFlag returns true if this user can flag.
func (u *User) CanFlag() bool {
	return u.Points > 50
}
