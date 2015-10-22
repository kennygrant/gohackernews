// Package comments represents the comment resource
package comments

import (
	"fmt"
	"strings"
	"time"

	"github.com/fragmenta/model"
	"github.com/fragmenta/model/validate"
	"github.com/fragmenta/query"

	"github.com/kennygrant/gohackernews/src/lib/status"
)

// Need rank on comments too? rank desc,
const RankOrder = "points desc, id desc"

// Comment handles saving and retreiving comments from the database
type Comment struct {
	model.Model
	status.ModelStatus
	Text string

	Points int64
	Rank   int64

	// Comment tree
	ParentId  int64
	DottedIds string
	Children  []*Comment

	//Join ids
	UserId  int64
	StoryId int64

	// Denormalised join details, for quick display, alternatively could use joins
	UserName  string
	StoryName string
}

// AllowedParams returns an array of allowed param keys
func AllowedParams() []string {
	return []string{"text"}
}

// AllowedParamsAdmin returns an array of allowed param keys
func AllowedParamsAdmin() []string {
	return []string{"status", "user_id", "user_name", "parent_id", "points", "story_id", "text", "dotted_ids", "story_id", "story_name", "rank"}
}

// NewWithColumns creates a new comment instance and fills it with data from the database cols provided
func NewWithColumns(cols map[string]interface{}) *Comment {

	comment := New()
	comment.Id = validate.Int(cols["id"])
	comment.CreatedAt = validate.Time(cols["created_at"])
	comment.UpdatedAt = validate.Time(cols["updated_at"])
	comment.Status = validate.Int(cols["status"])
	comment.Text = validate.String(cols["text"])
	comment.Points = validate.Int(cols["points"])
	comment.Rank = validate.Int(cols["rank"])

	comment.UserId = validate.Int(cols["user_id"])
	comment.UserName = validate.String(cols["user_name"])

	comment.StoryId = validate.Int(cols["story_id"])
	comment.StoryName = validate.String(cols["story_name"])

	comment.ParentId = validate.Int(cols["parent_id"])
	comment.DottedIds = validate.String(cols["dotted_ids"])

	return comment
}

// New creates and initialises a new comment instance
func New() *Comment {
	comment := &Comment{}
	comment.Model.Init()
	comment.Status = status.Draft
	comment.TableName = "comments"
	return comment
}

// Create inserts a new record in the database using params, and returns the newly created id
func Create(params map[string]string) (int64, error) {

	// Check params for invalid values
	err := validateParams(params)
	if err != nil {
		return 0, err
	}

	// Update date params
	params["created_at"] = query.TimeString(time.Now().UTC())
	params["updated_at"] = query.TimeString(time.Now().UTC())

	return Query().Insert(params)
}

// validateParams checks these params pass validation checks
func validateParams(params map[string]string) error {

	// Now check params are as we expect
	err := validate.Length(params["id"], 0, -1)
	if err != nil {
		return err
	}
	err = validate.Length(params["name"], 0, 255)
	if err != nil {
		return err
	}

	return err
}

// Find returns a single record by id in params
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
			if p.Id == c.ParentId {
				p.Children = append(p.Children, c)
				found = true
				break
			}
		}
		if !found {
			for _, p := range childComments {
				if p.Id == c.ParentId {
					p.Children = append(p.Children, c)
					break
				}
			}
		}
	}

	return rootComments, nil
}

// Query returns a new query for comments
func Query() *query.Query {
	p := New()
	return query.New(p.TableName, p.KeyName)
}

// Published returns a query for all comments with status >= published
func Published() *query.Query {
	return Query().Where("status>=?", status.Published)
}

// Where returns a Where query for comments with the arguments supplied
func Where(format string, args ...interface{}) *query.Query {
	return Query().Where(format, args...)
}

// Update sets the record in the database from params
func (m *Comment) Update(params map[string]string) error {

	// Check params for invalid values
	err := validateParams(params)
	if err != nil {
		return err
	}

	// Update date params
	params["updated_at"] = query.TimeString(time.Now().UTC())

	return Query().Where("id=?", m.Id).Update(params)
}

// NegativePoints returns a negative point score between 0 and 5 (positive points return 0, below -6 returns 6)
func (m *Comment) NegativePoints() int64 {
	if m.Points > 0 {
		return 0
	}
	if m.Points < -6 {
		return 6
	}

	return -m.Points
}

// Level returns the nesting level of this comment, based on dotted_ids
func (m *Comment) Level() int64 {
	if m.ParentId > 0 {
		return int64(strings.Count(m.DottedIds, "."))
	}
	return 0
}

// Root returns true if this is a root comment
func (m *Comment) Root() bool {
	return m.ParentId == 0
}

// Destroy removes the record from the database
func (m *Comment) Destroy() error {
	return Query().Where("id=?", m.Id).Delete()
}

// URLStory returns the internal resource URL for our story
func (m *Comment) URLStory() string {
	return fmt.Sprintf("/stories/%d", m.StoryId)
}
