package stats

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/fragmenta/router"
)

// Put this in a separate package called stats

// PurgeInterval is the interval at which users are purged from the current list
var PurgeInterval = time.Minute * 5

// identifiers holds a hash of anonymised user records
// obviously an in-memory store is not suitable for very large sites
// but for smaller sites with a few hundred concurrent users it's fine
var identifiers = make(map[string]time.Time)

// RegisterHit registers a hit and ups user count if required
func RegisterHit(context router.Context) {

	// Extract the IP from the address
	ip := context.Request().RemoteAddr
	forward := context.Request().Header.Get("X-Forwarded-For")
	if len(forward) > 0 {
		ip = forward
	}

	// Use UA as well as ip for unique values per browser session
	ua := context.Request().Header.Get("User-Agent")

	context.Logf("IP recd:%s ua:%s", ip, ua)

	// Hash for anonymity in our store
	hasher := sha1.New()
	hasher.Write([]byte(ip))
	hasher.Write([]byte(ua))
	id := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	// Insert the entry with current time
	identifiers[id] = time.Now()

}

// HandleUserCount serves a get request at /stats/users/count
func HandleUserCount(context router.Context) error {
	// Render json of our count for the javascript to display
	json := fmt.Sprintf("{\"users\":%d}", len(identifiers))
	_, err := context.Writer().Write([]byte(json))
	return err
}

// UserCount returns a count of users in the last 5 minutes
func UserCount() int {
	return len(identifiers)
}

// Clean up users list at intervals
func init() {
	purgeUsers()
}

// purgeUsers clears the users list of users who last acted PurgeInterval ago
func purgeUsers() {
	for k, v := range identifiers {
		purgeTime := time.Now().Add(-PurgeInterval)
		if v.Before(purgeTime) {
			delete(identifiers, k)
		}
	}

	//	fmt.Printf("Cleaning users:%d", len(identifiers))
	time.AfterFunc(time.Minute*1, purgeUsers)
}
