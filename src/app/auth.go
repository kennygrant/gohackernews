package app

import (
	"github.com/fragmenta/auth"
	"github.com/fragmenta/auth/can"
	"github.com/fragmenta/server/config"

	"github.com/kennygrant/gohackernews/src/comments"
	"github.com/kennygrant/gohackernews/src/stories"
	"github.com/kennygrant/gohackernews/src/users"
)

// SetupAuth sets up the auth pkg and authorisation for users
func SetupAuth() {

	// Set up the auth package with our secrets from config
	auth.HMACKey = auth.HexToBytes(config.Get("hmac_key"))
	auth.SecretKey = auth.HexToBytes(config.Get("secret_key"))
	auth.SessionName = config.Get("session_name")

	// Enable https cookies on production server - everyone should be on https
	if config.Production() {
		auth.SecureCookies = true
	}

	// Set up our authorisation for user roles on resources using can pkg

	// Admins are allowed to manage all resources
	can.Authorise(users.Admin, can.ManageResource, can.Anything)

	// Readers may edit their user
	can.AuthoriseOwner(users.Reader, can.UpdateResource, users.TableName)

	// Readers may add comments and edit their own comments
	can.Authorise(users.Reader, can.CreateResource, comments.TableName)
	can.AuthoriseOwner(users.Reader, can.UpdateResource, comments.TableName)

	// Readers may add stories and edit their own stories (up to time limit)
	can.Authorise(users.Reader, can.CreateResource, stories.TableName)
	can.AuthoriseOwner(users.Reader, can.UpdateResource, stories.TableName)

	// Anon may create users
	can.AuthoriseOwner(users.Anon, can.CreateResource, users.TableName)

}
