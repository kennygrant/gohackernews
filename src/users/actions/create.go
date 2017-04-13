package useractions

import (
	"fmt"
	"net/http"

	"github.com/fragmenta/auth"
	"github.com/fragmenta/auth/can"
	"github.com/fragmenta/mux"
	"github.com/fragmenta/server"
	"github.com/fragmenta/server/log"
	"github.com/fragmenta/view"

	"github.com/kennygrant/gohackernews/src/lib/session"
	"github.com/kennygrant/gohackernews/src/lib/status"
	"github.com/kennygrant/gohackernews/src/users"
)

// HandleCreateShow serves the create form via GET for users.
func HandleCreateShow(w http.ResponseWriter, r *http.Request) error {

	user := users.New()

	// Authorise
	err := can.Create(user, session.CurrentUser(w, r))
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Check they're not logged in already if so redirect.
	if !session.CurrentUser(w, r).Anon() {
		return server.Redirect(w, r, "/?warn=already_logged_in")
	}

	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.AddKey("user", user)
	view.AddKey("error", params.Get("error"))
	return view.Render()
}

// HandleCreate handles the POST of the create form for users
func HandleCreate(w http.ResponseWriter, r *http.Request) error {

	user := users.New()

	// Check the authenticity token
	err := session.CheckAuthenticity(w, r)
	if err != nil {
		return err
	}

	// Authorise
	err = can.Create(user, session.CurrentUser(w, r))
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Check they're not logged in already if so redirect.
	if !session.CurrentUser(w, r).Anon() {
		return server.Redirect(w, r, "/?warn=already_logged_in")
	}

	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Check a user doesn't exist with this name or email already
	name := params.Get("name")
	email := params.Get("email")
	pass := params.Get("password")

	// Name must be at least 2 characters
	if len(name) < 2 {
		return server.InternalError(err, "Name too short", "Sorry, names must be at least 2 characters long")
	}

	// Password must be at least 6 characters
	if len(pass) < 6 {
		return server.InternalError(err, "Password too short", "Sorry, passwords must be at least 6 characters long")
	}

	// Name is not optional so always check duplicates
	duplicates, err := users.FindAll(users.Where("name=?", name))
	if err != nil {
		return server.InternalError(err)
	}
	if len(duplicates) > 0 {
		return server.Redirect(w, r, "/users/create?error=duplicate_name")
	}

	// Email is optional, so allow blank email and don't check duplicates if so
	if email != "" {
		duplicates, err = users.FindAll(users.Where("email=?", email))
		if err != nil {
			return server.InternalError(err)
		}
		if len(duplicates) > 0 {
			return server.Redirect(w, r, "/users/create?error=duplicate_email")
		}
	}

	// Set the password hash from the password
	hash, err := auth.HashPassword(pass)
	if err != nil {
		return server.InternalError(err)
	}
	params.SetString("password_hash", hash)

	// Validate the params, removing any we don't accept
	userParams := user.ValidateParams(params.Map(), users.AllowedParams())

	// Set some defaults for the new user
	userParams["status"] = fmt.Sprintf("%d", status.Published)
	userParams["role"] = fmt.Sprintf("%d", users.Reader)
	userParams["points"] = "1"

	id, err := user.Create(userParams)
	if err != nil {
		return server.InternalError(err)
	}

	// Redirect to the new user
	user, err = users.Find(id)
	if err != nil {
		return server.InternalError(err)
	}

	// Log in automatically as the new user they have just created
	session, err := auth.Session(w, r)
	if err != nil {
		log.Info(log.V{"msg": "login failed", "email": user.Email, "user_id": user.ID, "status": http.StatusInternalServerError})
	}

	// Success, log it and set the cookie with user id
	session.Set(auth.SessionUserKey, fmt.Sprintf("%d", user.ID))
	session.Save(w)

	// Log action
	log.Info(log.V{"msg": "login success", "user_email": user.Email, "user_id": user.ID})

	return server.Redirect(w, r, "/")
}
