package mail

import (
	"fmt"

	"github.com/fragmenta/view"
	"github.com/sendgrid/sendgrid-go"
)

// The Mail service secret key/password (must be set before first sending)
var secret string

// The default sender
var from string

// Setup sets the user and secret for use in sending mail (possibly later we should have a config etc)
func Setup(s string, f string) {
	secret = s
	from = f
}

// Send sends mail
func Send(recipients []string, subject string, template string, context map[string]interface{}) error {

	// For now  ensure that we don't send to more than 1 recipient while we debug emails
	if len(recipients) > 1 {
		return fmt.Errorf("mail.send: #error bad recipients for debug %v", recipients)
	}

	if recipients[0] != "kennygrant@gmail.com" {
		return fmt.Errorf("mail.send: #error bad recipients for debug %v", recipients)
	}

	// Send via sendgrid
	sg := sendgrid.NewSendGridClientWithApiKey(secret)

	message := sendgrid.NewMail()
	message.SetFrom(from)
	message.AddTos(recipients)
	message.SetSubject(subject)

	// Load the template, and substitute using context
	// We should possibly set layout from caller too?
	view := view.NewWithPath("", nil)
	view.Template(template)
	view.Context(context)

	html, err := view.RenderToString()
	if err != nil {
		return err
	}
	message.SetHTML(html)

	// For debug, print message
	fmt.Printf("#info sending MAIL to:%s", recipients)

	return sg.Send(message)
}

// SendOne sends email to ONE recipient only
func SendOne(recipient string, subject string, template string, context map[string]interface{}) error {
	return Send([]string{recipient}, subject, template, context)
}
