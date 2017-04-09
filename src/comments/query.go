package comments

import (
	"time"

	"github.com/fragmenta/query"

	"github.com/kennygrant/gohackernews/src/lib/resource"
	"github.com/kennygrant/gohackernews/src/lib/status"
)

const (
	// TableName is the database table for this resource
	TableName = "comments"
	// KeyName is the primary key value for this resource
	KeyName = "id"
	// Order defines the default sort order in sql for this resource
	Order = "rank desc, points desc, id desc"
)

// AllowedParamsAdmin returns an array of allowed param keys for Update and Create.
func AllowedParamsAdmin() []string {
	return []string{"status", "dotted_ids", "parent_id", "points", "rank", "story_id", "story_name", "text", "user_id", "user_name"}
}

// AllowedParams returns an array of allowed param keys for Update and Create.
func AllowedParams() []string {
	return []string{"text"}
}

// NewWithColumns creates a new comment instance and fills it with data from the database cols provided.
func NewWithColumns(cols map[string]interface{}) *Comment {

	comment := New()
	comment.ID = resource.ValidateInt(cols["id"])
	comment.CreatedAt = resource.ValidateTime(cols["created_at"])
	comment.UpdatedAt = resource.ValidateTime(cols["updated_at"])
	comment.Status = resource.ValidateInt(cols["status"])
	comment.DottedIDs = resource.ValidateString(cols["dotted_ids"])
	comment.ParentID = resource.ValidateInt(cols["parent_id"])
	comment.Points = resource.ValidateInt(cols["points"])
	comment.Rank = resource.ValidateInt(cols["rank"])
	comment.StoryID = resource.ValidateInt(cols["story_id"])
	comment.StoryName = resource.ValidateString(cols["story_name"])
	comment.Text = resource.ValidateString(cols["text"])
	comment.UserID = resource.ValidateInt(cols["user_id"])
	comment.UserName = resource.ValidateString(cols["user_name"])

	return comment
}

// New creates and initialises a new comment instance.
func New() *Comment {
	comment := &Comment{}
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	comment.TableName = TableName
	comment.KeyName = KeyName
	comment.Status = status.Draft
	return comment
}

// FindFirst fetches a single comment record from the database using
// a where query with the format and args provided.
func FindFirst(format string, args ...interface{}) (*Comment, error) {
	result, err := Query().Where(format, args...).FirstResult()
	if err != nil {
		return nil, err
	}
	return NewWithColumns(result), nil
}

// Find fetches a single comment record from the database by id.
func Find(id int64) (*Comment, error) {
	result, err := Query().Where("id=?", id).FirstResult()
	if err != nil {
		return nil, err
	}
	return NewWithColumns(result), nil
}

// FindAll returns all results for this query
func FindAll(q *query.Query) ([]*Comment, error) {

	// Fetch query.Results from query
	results, err := q.Results()
	if err != nil {
		return nil, err
	}

	// Construct an array of comments constructed from the results
	// We do things a little differently, as we have a tree of comments
	// root comments are added to the list, others are held in another list
	// and added as children to rootComments

	var rootComments, childComments []*Comment
	for _, cols := range results {
		c := NewWithColumns(cols)
		if c.Root() {
			rootComments = append(rootComments, c)
		} else {
			childComments = append(childComments, c)
		}
	}

	// Now walk through child comments, assigning them to their parent

	// Walk through comments, adding those with no parent id to comments list
	// and others to the parent comment in root comments
	for _, c := range childComments {
		found := false
		for _, p := range rootComments {
			if p.ID == c.ParentID {
				p.Children = append(p.Children, c)
				found = true
				break
			}
		}
		if !found {
			for _, p := range childComments {
				if p.ID == c.ParentID {
					p.Children = append(p.Children, c)
					break
				}
			}
		}
	}

	return rootComments, nil
}

// Query returns a new query for comments with a default order.
func Query() *query.Query {
	return query.New(TableName, KeyName).Order(Order)
}

// Where returns a new query for comments with the format and arguments supplied.
func Where(format string, args ...interface{}) *query.Query {
	return Query().Where(format, args...)
}

// Published returns a query for all comments with status >= published.
func Published() *query.Query {
	return Query().Where("status>=?", status.Published)
}
