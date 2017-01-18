package session

import (
	"context"
	"net/http"
	"strconv"

	"github.com/fragmenta/auth"
	"github.com/fragmenta/mux"
	"github.com/fragmenta/server/log"

	"github.com/kennygrant/gohackernews/src/users"
)

// currentUserCtxKey is used as a key on context
// for saving the logged in user for the request.
var currentUserCtxKey = &ctxKey{}

type ctxKey struct{}

// CurrentUser returns the saved user (or an empty anon user)
// for the current session cookie
func CurrentUser(w http.ResponseWriter, r *http.Request) *users.User {

	// Extract the current user (if any) from context
	u := r.Context().Value(currentUserCtxKey)
	if u != nil {
		return u.(*users.User)
	}

	// Start with an anon user by default (role 0, id 0)
	user := &users.User{}

	// Build the session from the secure cookie, or create a new one
	session, err := auth.Session(w, r)
	if err != nil {
		log.Info(log.V{"msg": "session error", "error": err, "status": http.StatusInternalServerError})
		return user
	}

	// Fetch the current user record if we have one recorded in the session
	var id int64
	val := session.Get(auth.SessionUserKey)
	if len(val) > 0 {
		id, err = strconv.ParseInt(val, 10, 64)
		if err != nil {
			log.Info(log.V{"msg": "session error decoding", "val": val, "error": err, "status": http.StatusInternalServerError})
			return user
		}
	}

	if id > 0 {
		user, err = users.Find(id)
		if err != nil {
			log.Info(log.V{"msg": "session error user not found", "user_id": id, "error": err, "status": http.StatusNotFound})
			return user
		}
	}

	// If we have a user, save it to the context as current user
	// so that if we are called twice it is not expensive,
	// and so that views can have this key set automatically
	ctx := r.Context()
	ctx = context.WithValue(ctx, currentUserCtxKey, user)
	r.WithContext(ctx)

	return user
}

// clearSession clears the request session cookie entirely.
// If an error is encountered in processing params, the session is cleared.
func clearSession(w http.ResponseWriter, r *http.Request) error {
	// Clear the session
	session, err := auth.SessionGet(r)
	if err != nil {
		return err
	}
	session.Clear(w)
	return nil
}

// CheckAuthenticity checks the authenticity token in params against cookie -
// The masked token is inserted into forms and POSTS by js.
// The token is inserted into the cookie by the middleware above.
// This is a shortcut for where you don't need params otherwise.
func CheckAuthenticity(w http.ResponseWriter, r *http.Request) error {

	// We should never be called on GET requests
	if r.Method == http.MethodGet {
		return nil
	}

	// Get the token from params and compare against cookie
	params, err := mux.Params(r)
	if err != nil {
		clearSession(w, r)
		return err
	}

	// Get the token from params (it is inserted there by js)
	// we do this to allow just one token in the head of every page
	token := params.Get(auth.SessionTokenKey)

	// Compare that param against the token stored in the session cookie.
	err = auth.CheckAuthenticityToken(token, r)
	if err != nil {
		clearSession(w, r)
		return err
	}

	return nil
}
