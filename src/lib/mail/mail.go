// Package mail provides a wrapper around sending mail via the fragile sendgrid API
package mail

import (
	"fmt"

	"github.com/fragmenta/view"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// The Mail service secret key/password (must be set before first sending)
var secret string

// The default sender
var from string

// Setup sets the user and secret for use in sending mail
func Setup(s string, f string) {
	secret = s
	from = f
}

// Send sends mail (using sendgrid API v3)
func Send(recipients []string, subject string, template string, context map[string]interface{}) error {

	// For now  ensure that we don't send to more than 1 recipient while we debug emails
	if len(recipients) > 1 {
		return fmt.Errorf("mail.send: #error bad recipients for debug %v", recipients)
	}

	if recipients[0] != "kennygrant@gmail.com" {
		return fmt.Errorf("mail.send: #error bad recipients for debug %v", recipients)
	}

	// Send via sendgrid
	// Apparently this API will probably break again without warning !
	// consider vendoring

	// Load the template, and substitute using context
	// We should possibly set layout from caller too?
	view := view.NewWithPath("", nil)
	view.Template(template)
	view.Context(context)

	html, err := view.RenderToString()
	if err != nil {
		return err
	}

	// Create a sendgrid message with v3
	sendgridContent := mail.NewContent("text/html", html)
	var sendgridRecipients []*mail.Email
	for _, r := range recipients {
		sendgridRecipients = append(sendgridRecipients, mail.NewEmail("", r))
	}

	message := mail.NewV3Mail()
	message.Subject = subject
	message.From = mail.NewEmail("", from)
	p := mail.NewPersonalization()
	p.AddTos(sendgridRecipients...)
	message.AddPersonalizations(p)
	message.AddContent(sendgridContent)

	request := sendgrid.GetRequest(secret, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(message)
	_, err = sendgrid.API(request)

	// For debug, print message
	fmt.Printf("#info sending MAIL to:%s", recipients)

	return err
}

// SendOne sends email to ONE recipient only
func SendOne(recipient string, subject string, template string, context map[string]interface{}) error {
	return Send([]string{recipient}, subject, template, context)
}
