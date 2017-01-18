package helpers

import (
	"strings"
	"testing"

	"github.com/fragmenta/server/config"
)

// TestRootURL tests our root url loaded
func TestRootURL(t *testing.T) {

	// Load the appropriate config - fix path here
	c := config.New()
	err := c.Load("./../../../secrets/fragmenta.json")
	if err != nil {
		t.Errorf("helpers: failed to load config")
	}
	config.Current = c
	config.Current.Mode = config.ModeTest

	if len(RootURL()) == 0 {
		t.Errorf("helpers: failed to load root url")
	}

}

// TestRootURL tests our root url loaded
func TestMarkup(t *testing.T) {

	// Use a table test here instead
	test := " @kenny "
	got := Markup(test)
	expected := "<a href="
	if !strings.Contains(string(got), expected) {
		t.Errorf("helpers: failed to convert markup expected:%s got:%s", expected, got)
	}

}
