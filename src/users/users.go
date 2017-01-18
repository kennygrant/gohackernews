// Package users represents the user resource
package users

import (
	"time"

	"github.com/kennygrant/gohackernews/src/lib/resource"
	"github.com/kennygrant/gohackernews/src/lib/status"
)

// User handles saving and retreiving users from the database
type User struct {
	// resource.Base defines behaviour and fields shared between all resources
	resource.Base

	// status.ResourceStatus defines a status field and associated behaviour
	status.ResourceStatus

	Email   string
	Name    string
	Points  int64
	Role    int64
	Summary string
	Text    string
	Title   string

	PasswordHash    string
	PasswordResetAt time.Time
}
