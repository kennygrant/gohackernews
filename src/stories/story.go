// Package stories represents the story resource
package stories

import (
	"fmt"
	"strings"
	"time"

	"github.com/fragmenta/model/file"

	"github.com/kennygrant/gohackernews/src/lib/resource"
	"github.com/kennygrant/gohackernews/src/lib/status"
)

// Story handles saving and retreiving stories from the database
type Story struct {
	// resource.Base defines behaviour and fields shared between all resources
	resource.Base

	// status.ResourceStatus defines a status field and associated behaviour
	status.ResourceStatus

	CommentCount int64
	Name         string
	Summary      string
	URL          string
	Points       int64
	Rank         int64
	UserID       int64

	// UserName denormalises the user name - use join instead?
	UserName string
}

// NegativePoints returns a negative point score or 0 if points is above 0
func (s *Story) NegativePoints() int64 {
	if s.Points > 0 {
		return 0
	}
	return -s.Points
}

// Domain returns the domain of the story URL
func (s *Story) Domain() string {
	parts := strings.Split(s.URL, "/")
	if len(parts) > 2 {
		return strings.Replace(parts[2], "www.", "", 1)
	}

	if len(s.URL) > 0 {
		return s.URL
	}

	return "GN"
}

// ShowAsk returns true if this is a Show: or Ask: story
func (s *Story) ShowAsk() bool {
	return strings.HasPrefix(s.Name, "Show:") || strings.HasPrefix(s.Name, "Ask:")
}

// DestinationURL returns the URL of the story
// if no url is set, it uses the CanonicalURL
func (s *Story) DestinationURL() string {
	if s.URL != "" {
		return s.URL
	}
	return s.CanonicalURL()
}

// PrimaryURL returns the URL to use for this story in lists
// Videos and Show Ask stories link to the story, for other links for now it is the destination
func (s *Story) PrimaryURL() string {
	// If video or show or ask, return story url
	if s.YouTube() || s.ShowAsk() {
		return s.CanonicalURL()
	}

	// If no url, return canonical url
	if s.URL == "" {
		return s.CanonicalURL()
	}

	return s.URL
}

// CanonicalURL is the canonical URL of the story on this site
// including a slug for seo
func (s *Story) CanonicalURL() string {
	return fmt.Sprintf("/stories/%d-%s", s.ID, file.SanitizeName(s.Name))
}

// Code returns true if this is a link to a git repository
// At present we only check for github urls, we should at least check for bitbucket
func (s *Story) Code() bool {
	if strings.Contains(s.URL, "https://github.com") {
		if strings.Contains(s.URL, "/commit/") || strings.Contains(s.URL, "/releases/") || strings.HasSuffix(s.URL, ".md") {
			return false
		}
		return true
	}
	return false
}

// GodocURL returns the godoc.org URL for this story, or empty string if none
func (s *Story) GodocURL() string {
	if s.Code() {
		return strings.Replace(s.URL, "https://github.com", "https://godoc.org/github.com", 1)
	}
	return ""
}

// VetURL returns a URL for goreportcard.com, for code repos
func (s *Story) VetURL() string {
	if s.Code() {
		return strings.Replace(s.URL, "https://github.com/", "http://goreportcard.com/report/", 1)
	}
	return ""
}

// YouTube returns true if this is a youtube video
func (s *Story) YouTube() bool {
	return strings.Contains(s.URL, "youtube.com/watch?v=")
}

// YouTubeURL returns the youtube URL
func (s *Story) YouTubeURL() string {
	url := strings.Replace(s.URL, "https://s.youtube.com", "https://www.youtube.com", 1)
	// https://www.youtube.com/watch?v=sZx3oZt7LVg ->
	// https://www.youtube.com/embed/sZx3oZt7LVg
	url = strings.Replace(url, "watch?v=", "embed/", 1)
	return url
}

// CommentCountDisplay returns the comment count or ellipsis if count is 0
func (s *Story) CommentCountDisplay() string {
	if s.CommentCount > 0 {
		return fmt.Sprintf("%d", s.CommentCount)
	}
	return "…"
}

// NameDisplay returns a title string without hashtags (assumed to be at the end),
// by truncating the title at the first #
func (s *Story) NameDisplay() string {
	if strings.Contains(s.Name, "#") {
		return s.Name[0:strings.Index(s.Name, "#")]
	}
	return s.Name
}

// Tags are defined as words beginning with # in the title
// TODO: for speed and clarity we could extract at submit time instead and store in db
func (s *Story) Tags() []string {
	var tags []string
	if strings.Contains(s.Name, "#") {
		// split on " #"
		parts := strings.Split(s.Name, " #")
		tags = parts[1:]
	}
	return tags
}

// Editable returns true if this story is editable.
// Stories are editable if less than 1 hours old
func (s *Story) Editable() bool {
	return time.Now().Sub(s.CreatedAt) < time.Hour*1
}

// OwnedBy returns true if this user id owns this story.
func (s *Story) OwnedBy(uid int64) bool {
	return uid == s.UserID
}
