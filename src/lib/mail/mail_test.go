package mail

import (
	"testing"

	"github.com/fragmenta/view"

	"github.com/kennygrant/gohackernews/src/lib/helpers"
)

// TestMail tests that mail formats properly in dev mode
func TestMail(t *testing.T) {

	view.Helpers["markup"] = helpers.Markup
	view.Helpers["timeago"] = helpers.TimeAgo
	view.Helpers["root_url"] = helpers.RootURL

	// In order to test, we rely on the view pkg being set up
	err := view.LoadTemplatesAtPaths([]string{"../.."}, view.Helpers)
	if err != nil {
		t.Errorf("mail: failed to load views:%s", err)
	}

	context := Context{
		"msg": "hello world",
	}

	recipient := "recipient@example.com"
	email := New(recipient)
	email.ReplyTo = "sender@example.com"
	email.Subject = "sub"
	email.Body = "<h1>Body</h1>"
	Send(email, context)

	// Try render
	email.Body, err = RenderTemplate(email, context)
	if err != nil {
		t.Errorf("mail: failed to render message :%s", err)
	}

}
