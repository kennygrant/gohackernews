// Package stories represents the story resource
package stories

import (
	"fmt"
	"strings"
	"time"

	"github.com/fragmenta/model"
	"github.com/fragmenta/model/validate"
	"github.com/fragmenta/query"
	"github.com/fragmenta/router"

	"github.com/kennygrant/gohackernews/src/lib/status"
)

// Story handles saving and retreiving stories from the database
type Story struct {
	model.Model
	status.ModelStatus
	Name    string
	Summary string
	Url     string
	Rank    int64
	Points  int64

	UserId       int64
	UserName     string
	CommentCount int64
}

// AllowedParams returns an array of allowed param keys
func AllowedParams() []string {
	return []string{"name", "points", "summary", "url"}
}

// AllowedParamsAdmin returns an array of allowed param keys
func AllowedParamsAdmin() []string {
	return []string{"status", "user_id", "user_name", "name", "points", "summary", "url", "comment_count", "tweeted_at", "newsletter_at"}
}

// NewWithColumns creates a new story instance and fills it with data from the database cols provided
func NewWithColumns(cols map[string]interface{}) *Story {

	story := New()
	story.Id = validate.Int(cols["id"])
	story.CreatedAt = validate.Time(cols["created_at"])
	story.UpdatedAt = validate.Time(cols["updated_at"])
	story.Status = validate.Int(cols["status"])
	story.Name = validate.String(cols["name"])
	story.Summary = validate.String(cols["summary"])
	story.Url = validate.String(cols["url"])
	story.Rank = validate.Int(cols["rank"])
	story.Points = validate.Int(cols["points"])

	story.UserId = validate.Int(cols["user_id"])
	story.UserName = validate.String(cols["user_name"])
	story.CommentCount = validate.Int(cols["comment_count"])

	return story
}

// New creates and initialises a new story instance
func New() *Story {
	story := &Story{}
	story.Model.Init()
	story.Status = status.Draft
	story.TableName = "stories"
	return story
}

// Create inserts a new record in the database using params, and returns the newly created id
func Create(params map[string]string) (int64, error) {

	// Check params for invalid values
	err := validateParams(params, true)
	if err != nil {
		return 0, err
	}

	// Update date params
	params["created_at"] = query.TimeString(time.Now().UTC())
	params["updated_at"] = query.TimeString(time.Now().UTC())

	return Query().Insert(params)
}

// validateParams checks these params pass validation checks
// TODO: reconsider best interface for this - don't like the bool
func validateParams(params map[string]string, checkAll bool) error {

	if checkAll || len(params["name"]) > 0 {
		err := validate.Length(params["name"], 2, 300)
		if err != nil {
			return router.BadRequestError(err, "Invalid Name", "The name must be over 2 characters")
		}
	}

	if checkAll || len(params["url"]) > 0 {

		err := validate.Length(params["url"], 5, 1000)
		if err != nil {
			return router.BadRequestError(err, "Invalid URL", "The URL must be over 5 characters")
		}

		if !strings.HasPrefix(params["url"], "http://") && !strings.HasPrefix(params["url"], "https://") {
			return router.BadRequestError(nil, "Invalid URL", "The URL must have scheme https:// or http://")
		}

	}

	return nil
}

// Find returns a single record by id in params
func Find(id int64) (*Story, error) {
	result, err := Query().Where("id=?", id).FirstResult()
	if err != nil {
		return nil, err
	}
	return NewWithColumns(result), nil
}

// FindAll returns all results for this query
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

// Query returns a new query for stories
func Query() *query.Query {
	p := New()
	return query.New(p.TableName, p.KeyName)
}

// Where returns a Where query for stories with the arguments supplied
func Where(format string, args ...interface{}) *query.Query {
	return Query().Where(format, args...)
}

// Published returns a query for all stories with status >= published
func Published() *query.Query {
	return Query().Where("status>=?", status.Published)
}

// Popular returns a query for all stories with points over a certain threshold
func Popular() *query.Query {
	return Query().Where("points > 2")
}

// Update sets the record in the database from params
func (m *Story) Update(params map[string]string) error {

	// Check params for invalid values, but only if passed in
	err := validateParams(params, false)
	if err != nil {
		return err
	}

	// Update date params
	params["updated_at"] = query.TimeString(time.Now().UTC())

	return Query().Where("id=?", m.Id).Update(params)
}

// Destroy removes the record from the database
func (m *Story) Destroy() error {
	return Query().Where("id=?", m.Id).Delete()
}

// NegativePoints returns a negative point score or 0 if points is above 0
func (m *Story) NegativePoints() int64 {
	if m.Points > 0 {
		return 0
	}
	return -m.Points
}

// Domain returns the domain of the story URL
func (m *Story) Domain() string {
	parts := strings.Split(m.Url, "/")
	if len(parts) > 2 {
		return strings.Replace(parts[2], "www.", "", 1)
	}

	if len(m.Url) > 0 {
		return m.Url
	}

	return "GN"
}

// ShowAsk returns true if this is a Show: or Ask: story
func (m *Story) ShowAsk() bool {
	return strings.HasPrefix(m.Name, "Show:") || strings.HasPrefix(m.Name, "Ask:")
}

// DestinationURL returns the URL of the story (either set URL or if unset for Ask:, just the ShowURL)
func (m *Story) DestinationURL() string {
	if m.Url != "" {
		return m.Url
	}
	return m.URLShow()
}

// ListURL returns the URL to use for this story in lists for Show/Ask stories, this is the story link, for others, it is the destination URL
func (m *Story) ListURL() string {
	if !m.ShowAsk() && m.Url != "" {
		return m.Url
	}
	return m.URLShow()
}

// Code returns true if this is a link to a git repository
// At present we only check for github urls, we should at least check for bitbucket
func (m *Story) Code() bool {
	if strings.Contains(m.Url, "https://github.com") {
		if strings.Contains(m.Url, "/commit/") || strings.HasSuffix(m.Url, ".md") {
			return false
		}
		return true
	}
	return false
}

// GodocURL returns the godoc.org URL for this story, or empty string if none
func (m *Story) GodocURL() string {
	if m.Code() {
		return strings.Replace(m.Url, "https://github.com", "https://godoc.org/github.com", 1)
	}
	return ""
}

// VetURL returns a URL for goreportcard.com, for code repos
func (m *Story) VetURL() string {
	if m.Code() {
		return strings.Replace(m.Url, "https://github.com/", "http://goreportcard.com/report/", 1)
	}
	return ""
}

// YouTube returns true if this is a youtube video
func (m *Story) YouTube() bool {
	return strings.Contains(m.Url, "youtube.com/watch?v=")
}

// YouTubeURL returns the youtube URL
func (m *Story) YouTubeURL() string {
	url := strings.Replace(m.Url, "https://m.youtube.com", "https://www.youtube.com", 1)
	// https://www.youtube.com/watch?v=sZx3oZt7LVg ->
	// https://www.youtube.com/embed/sZx3oZt7LVg
	url = strings.Replace(url, "watch?v=", "embed/", 1)
	return url
}

// CommentCountDisplay returns the comment count or ellipsis if count is 0
func (m *Story) CommentCountDisplay() string {
	if m.CommentCount > 0 {
		return fmt.Sprintf("%d", m.CommentCount)
	}
	return "…"
}

// NameDisplay returns a title string without hashtags (assumed to be at the end),
// by truncating the title at the first #
func (m *Story) NameDisplay() string {
	if strings.Contains(m.Name, "#") {
		return m.Name[0:strings.Index(m.Name, "#")]
	}
	return m.Name
}

// Tags are defined as words beginning with # in the title
// TODO: for speed and clarity we could extract at submit time instead and store in db
func (m *Story) Tags() []string {
	var tags []string
	if strings.Contains(m.Name, "#") {
		// split on " #"
		parts := strings.Split(m.Name, " #")
		tags = parts[1:]
	}
	return tags
}
