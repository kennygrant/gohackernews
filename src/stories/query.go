package stories

import (
	"time"

	"github.com/fragmenta/query"

	"github.com/kennygrant/gohackernews/src/lib/resource"
	"github.com/kennygrant/gohackernews/src/lib/status"
)

const (
	// TableName is the database table for this resource
	TableName = "stories"
	// KeyName is the primary key value for this resource
	KeyName = "id"
	// Order defines the default sort order in sql for this resource
	Order = "name asc, id desc"
)

// AllowedParams returns the cols editable by everyone
func AllowedParams() []string {
	return []string{"name", "summary", "url"}
}

// AllowedParamsAdmin returns the cols editable by admins
func AllowedParamsAdmin() []string {
	return []string{"status", "comment_count", "name", "points", "rank", "summary", "url", "user_id", "user_name"}
}

// NewWithColumns creates a new story instance and fills it with data from the database cols provided.
func NewWithColumns(cols map[string]interface{}) *Story {

	story := New()
	story.ID = resource.ValidateInt(cols["id"])
	story.CreatedAt = resource.ValidateTime(cols["created_at"])
	story.UpdatedAt = resource.ValidateTime(cols["updated_at"])
	story.Status = resource.ValidateInt(cols["status"])
	story.CommentCount = resource.ValidateInt(cols["comment_count"])
	story.Name = resource.ValidateString(cols["name"])
	story.Points = resource.ValidateInt(cols["points"])
	story.Rank = resource.ValidateInt(cols["rank"])
	story.Summary = resource.ValidateString(cols["summary"])
	story.URL = resource.ValidateString(cols["url"])
	story.UserID = resource.ValidateInt(cols["user_id"])
	story.UserName = resource.ValidateString(cols["user_name"])

	return story
}

// New creates and initialises a new story instance.
func New() *Story {
	story := &Story{}
	story.CreatedAt = time.Now()
	story.UpdatedAt = time.Now()
	story.TableName = TableName
	story.KeyName = KeyName
	story.Status = status.Draft
	return story
}

// FindFirst fetches a single story record from the database using
// a where query with the format and args provided.
func FindFirst(format string, args ...interface{}) (*Story, error) {
	result, err := Query().Where(format, args...).FirstResult()
	if err != nil {
		return nil, err
	}
	return NewWithColumns(result), nil
}

// Find fetches a single story record from the database by id.
func Find(id int64) (*Story, error) {
	result, err := Query().Where("id=?", id).FirstResult()
	if err != nil {
		return nil, err
	}
	return NewWithColumns(result), nil
}

// FindAll fetches all story records matching this query from the database.
func FindAll(q *query.Query) ([]*Story, error) {

	// Fetch query.Results from query
	results, err := q.Results()
	if err != nil {
		return nil, err
	}

	// Return an array of stories constructed from the results
	var stories []*Story
	for _, cols := range results {
		p := NewWithColumns(cols)
		stories = append(stories, p)
	}

	return stories, nil
}

// Popular returns a query for all stories with points over a certain threshold
func Popular() *query.Query {
	return Query().Where("points > 2")
}

// Query returns a new query for stories with a default order.
func Query() *query.Query {
	return query.New(TableName, KeyName).Order(Order)
}

// Where returns a new query for stories with the format and arguments supplied.
func Where(format string, args ...interface{}) *query.Query {
	return Query().Where(format, args...)
}

// Published returns a query for all stories with status >= published.
func Published() *query.Query {
	return Query().Where("status>=?", status.Published)
}
