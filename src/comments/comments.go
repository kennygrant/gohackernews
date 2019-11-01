// Package comments represents the comment resource
package comments

import (
	"fmt"
	"strings"
	"time"

	"github.com/kennygrant/gohackernews/src/lib/resource"
	"github.com/kennygrant/gohackernews/src/lib/status"
)

// Comment handles saving and retreiving comments from the database
type Comment struct {
	// resource.Base defines behaviour and fields shared between all resources
	resource.Base

	// status.ResourceStatus defines a status field and associated behaviour
	status.ResourceStatus

	// A path to the story through parents
	DottedIDs string
	// The parent comment id (may be 0)
	ParentID int64
	// Any child comments (may be empty)
	Children []*Comment

	// Score of the comment (raw points and calculated rank with gravity)
	Points int64
	Rank   int64

	// The main story text
	Text string

	// Joined story and user
	StoryID int64
	UserID  int64

	// Denormalised attributes from joins
	StoryName string
	UserName  string
}

// Level returns the nesting level of this comment, based on dotted_ids
func (c *Comment) Level() int64 {
	if c.ParentID > 0 {
		return int64(strings.Count(c.DottedIDs, "."))
	}
	return 0
}

// Root returns true if this is a root comment
func (c *Comment) Root() bool {
	return c.ParentID == 0
}

// Destroy removes the record from the database
func (c *Comment) Destroy() error {
	return Query().Order("").Where("id=?", c.ID).Delete()
}

// StoryURL returns the internal resource URL for our story
func (c *Comment) StoryURL() string {
	return fmt.Sprintf("/stories/%d", c.StoryID)
}

// Editable returns true if this comment is editable.
// Comments are editable if less than 3 hours old.
func (c *Comment) Editable() bool {
	return time.Now().Sub(c.CreatedAt) < time.Hour*3
}

// OwnedBy returns true if this user id owns this comment.
func (c *Comment) OwnedBy(uid int64) bool {
	return uid == c.UserID
}

// NegativePoints returns a negative point score between 0 and 5 (positive points return 0, below -6 returns 6)
func (c *Comment) NegativePoints() int64 {
	if c.Points > 0 {
		return 0
	}
	if c.Points < -6 {
		return 6
	}

	return -c.Points
}
