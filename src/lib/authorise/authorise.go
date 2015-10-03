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

	// If not public path, check based on user role
	user := c.Get("current_user").(*users.User)
	switch user.Role {
	case users.RoleAdmin:
		return authoriseAdmin(c, r)
	default:
		return authoriseReader(c, r)
	}

}

// ResourceAndAuthenticity authorises the path and resource for the current user
func ResourceAndAuthenticity(c router.Context, r ResourceModel) error {

	// Check the authenticity token first
	err := AuthenticityToken(c)
	if err != nil {
		return err
	}

	// Now authorise the resource as normal
	return Resource(c, r)
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

	// Allow upvotes and downvotes
	if strings.HasSuffix(c.Path(), "/upvote") && user.CanUpvote() {
		return nil
	}

	if strings.HasSuffix(c.Path(), "/downvote") && user.CanDownvote() {
		return nil
	}

	if r != nil {
		if r.OwnedBy(user.Id) {
			return nil
		}
	}

	return fmt.Errorf("Path and Resource not authorized:%s %v", c.Path(), r)

}
