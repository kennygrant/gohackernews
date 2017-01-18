// Package helpers contains view helpers
package helpers

import (
	"fmt"
	"html/template"
	"time"

	"github.com/fragmenta/server/config"
	"github.com/fragmenta/view/helpers"

	"github.com/kennygrant/gohackernews/src/lib/text"
)

// RootURL returns the root url
func RootURL() string {
	return config.Get("root_url")
}

// Markup converts text from stories into sanitized html
func Markup(s string) template.HTML {

	// Convert bare links and usernames to anchors
	s = text.ConvertLinks(s)

	// Convert newlimnes to paragraph tags
	s = text.ConvertNewlines(s)

	// Run sanitize on the resulting string
	// (parses html and removes unwated tags and attributes)
	return helpers.Sanitize(s)
}

// TimeAgo returns a string for a time in format x seconds ago
func TimeAgo(d time.Time) string {

	duration := time.Since(d)
	hours := duration / time.Hour

	switch {
	case duration < time.Minute:
		return fmt.Sprintf("%d seconds ago", duration/time.Second)
	case duration < time.Hour:
		return fmt.Sprintf("%d minutes ago", duration/time.Minute)
	case duration < time.Hour*24:
		unit := "hour"
		if hours > 1 {
			unit = "hours"
		}
		return fmt.Sprintf("%d %s ago", hours, unit)
	default:
		unit := "day"
		if hours > 48 {
			unit = "days"
		}
		return fmt.Sprintf("%d %s ago", hours/24, unit)
	}

}
