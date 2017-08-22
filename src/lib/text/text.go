// Package text performs text manipulation on html strings
// it is used by stories and comments
package text

import (
	"regexp"
	"strings"
)

var (
	// Trailing defines optional characters allowed after a url or username
	// this excludes some valid urls but urls are not expected to end
	// in these characters
	trailing = `([\s!?.,]?)`

	// Search for links prefaced by word separators (\s\n\.)
	// i.e. not already in anchors, and replace with auto-link
	//	`\s(https?://.*)[\s.!?]`
	// match urls at start of text or with space before only
	urlRx = regexp.MustCompile(`(\A|[\s]+)(https?://[^\s><]*)` + trailing)

	// Search for \s@name in text and replace with links to username search
	// requires an endpoint that redirects /u/kenny to /users/1 etc.
	userRx = regexp.MustCompile(`(\A|[\s]+)@([^\s!?.,<>]*)` + trailing)

	// Search for trailing <p>\s for ConvertNewlines
	trailingPara = regexp.MustCompile(`<p>\s*\z`)
)

// ConvertNewlines converts \n to paragraph tags
// if the text already contains paragraph tags, return unaltered
func ConvertNewlines(s string) string {
	if strings.Contains(s, "<p>") {
		return s
	}

	// Start with para
	s = "<p>" + s
	// Replace newlines with paras
	s = strings.Replace(s, "\n", "</p><p>", -1)
	// Remove trailing para added in step above
	s = string(trailingPara.ReplaceAll([]byte(s), []byte("")))
	return s
}

// ConvertLinks returns the text with various transformations applied -
// bare links are turned into anchor tags, and @refs are turned into user links.
// this is somewhat fragile, better to parse the html
func ConvertLinks(s string) string {
	bytes := []byte(s)
	// Replace bare links with active links
	bytes = urlRx.ReplaceAll(bytes, []byte(`$1<a href="$2">$2</a>$3`))
	// Replace usernames with links
	bytes = userRx.ReplaceAll(bytes, []byte(`$1<a href="/u/$2">@$2</a>$3`))
	return string(bytes)
}
