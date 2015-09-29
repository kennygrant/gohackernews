package useractions

import (
	"fmt"

	"github.com/fragmenta/auth"
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/kennygrant/hackernews/src/lib/authorise"
	"github.com/kennygrant/hackernews/src/users"
)

// HandleLoginShow shows the page at /users/login
func HandleLoginShow(context router.Context) error {
	// Setup context for template
	view := view.New(context)

	// Check we're not already logged in, if so redirect with a message
	// we could alternatively display an error here?
	if !authorise.CurrentUser(context).Anon() {
		return router.Redirect(context, "/?warn=already_logged_in")
	}

	switch context.Param("error") {
	case "failed_email":
		view.AddKey("warning", "Sorry, we couldn't find a user with that email.")
	case "failed_password":
		view.AddKey("warning", "Sorry, the password was incorrect, please try again.")
	}

	// Serve
	return view.Render()
}

// HandleLogin handles a post to /users/login
func HandleLogin(context router.Context) error {

	// Check we're not already logged in, if so redirect

	// Get the user details from the database
	params, err := context.Params()
	if err != nil {
		return router.NotFoundError(err)
	}

	// Need something neater than this - how best to do it?
	q := users.Where("email=?", params.Get("email"))
	user, err := users.First(q)
	if err != nil {
		context.Logf("#error Login failed for user no such user : %s %s", params.Get("email"), err)
		return router.Redirect(context, "/users/login?error=failed_email")

	}

	err = auth.CheckPassword(params.Get("password"), user.EncryptedPassword)

	if err != nil {
		context.Logf("#error Login failed for user : %s %s", params.Get("email"), err)
		return router.Redirect(context, "/users/login?error=failed_password")
	}

	// Save the fact user is logged in to session cookie
	err = loginUser(context, user)
	if err != nil {
		return router.InternalError(err)
	}

	// Redirect to whatever page the user tried to visit before (if any)
	// For now send them to root
	return router.Redirect(context, "/")

}

func loginUser(context router.Context, user *users.User) error {
	// Now save the user details in a secure cookie, so that we remember the next request
	session, err := auth.Session(context, context.Request())
	if err != nil {
		return err
	}

	context.Logf("#info Login success for user: %d %s", user.Id, user.Email)
	session.Set(auth.SessionUserKey, fmt.Sprintf("%d", user.Id))
	session.Save(context)
	return nil
}
