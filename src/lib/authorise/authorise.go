package authorise

import (
	"fmt"
	"strings"

	"github.com/fragmenta/auth"
	"github.com/fragmenta/router"
	"github.com/fragmenta/server"

	"github.com/kennygrant/hackernews/src/users"
)

// ResourceModel defines the interface for models passed to authorise.Resource
type ResourceModel interface {
	OwnedBy(int64) bool
}

// Setup authentication and authorization keys for this app
func Setup(s *server.Server) {

	// Set up our secret keys which we take from the config
	// NB these are hex strings which we convert to bytes, for ease of presentation in secrets file
	c := s.Configuration()
	auth.HMACKey = auth.HexToBytes(c["hmac_key"])
	auth.SecretKey = auth.HexToBytes(c["secret_key"])
	auth.SessionName = "gohackernews"

	// Enable https cookies on production server - we don't have https, so don't do this
	//	if s.Production() {
	//		auth.SecureCookies = true
	//	}

}

// CurrentUserFilter returns a filter function which sets the current user on the context
func CurrentUserFilter(c router.Context) error {
	u := CurrentUser(c)
	c.Set("current_user", u)
	return nil
}

// Path authorises the path for the current user
func Path(c router.Context) error {
	return Resource(c, nil)
}

// Resource authorises the path and resource for the current user
// if model is nil it is ignored and permission granted
func Resource(c router.Context, r ResourceModel) error {

	// Short circuit evaluation if this is a public path
	if publicPath(c.Path()) {
		return nil
	}

	// If not public path, check based on user role
	user := c.Get("current_user").(*users.User)
	switch user.Role {
	case users.RoleAdmin:
		return authoriseAdmin(c, r)
	default:
		return authoriseReader(c, r)
	}

}

// Admins can see all screens
func authoriseAdmin(c router.Context, r ResourceModel) error {
	return nil
}

// authoriseReader returns error if the path/resource is not authorised
func authoriseReader(c router.Context, r ResourceModel) error {
	user := c.Get("current_user").(*users.User)

	if c.Path() == "/stories/create" && user.CanSubmit() {
		return nil
	}

	if c.Path() == "/comments/create" && user.CanComment() {
		return nil
	}

	if r != nil {
		if r.OwnedBy(user.Id) {
			return nil
		}
	}

	return fmt.Errorf("Path and Resource not authorized:%s %v", c.Path(), r)

}

// publicPath returns true if this path should always be allowed, regardless of user role
func publicPath(p string) bool {
	if p == "/" {
		return true
	}

	// Anon can log in and create a user (register)
	if p == "/users/login" || p == "/users/create" {
		return true
	}

	// This is a way of saying index and show only are public, no actions
	// TODO: find a neater way to do this?
	if strings.HasPrefix(p, "/comments") && strings.Count(p, "/") < 3 {
		return true
	}
	if strings.HasPrefix(p, "/stories") && strings.Count(p, "/") < 3 {
		return true
	}

	return false
}
