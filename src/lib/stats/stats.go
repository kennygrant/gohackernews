package stats

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Put this in a separate package called stats

// PurgeInterval is the interval at which users are purged from the current list
var PurgeInterval = time.Minute * 5

// identifiers holds a hash of anonymised user records
// obviously an in-memory store is not suitable for very large sites
// but for smaller sites with a few hundred concurrent users it's fine
var identifiers = make(map[string]time.Time)
var mu sync.RWMutex

// RegisterHit registers a hit and ups user count if required
func RegisterHit(r *http.Request) {

	// Use UA as well as ip for unique values per browser session
	ua := r.Header.Get("User-Agent")
	// Ignore obvious bots (Googlebot etc)
	if strings.Contains(ua, "bot") {
		return
	}
	// Ignore requests for xml (assumed to be feeds or sitemap)
	if strings.HasSuffix(r.URL.Path, ".xml") {
		return
	}

	// Extract the IP from the address
	ip := r.RemoteAddr
	forward := r.Header.Get("X-Forwarded-For")
	if len(forward) > 0 {
		ip = forward
	}

	// Hash for anonymity in our store
	hasher := sha1.New()
	hasher.Write([]byte(ip))
	hasher.Write([]byte(ua))
	id := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	// Insert the entry with current time
	mu.Lock()
	identifiers[id] = time.Now()
	mu.Unlock()
}

// HandleUserCount serves a get request at /stats/users/count
func HandleUserCount(w http.ResponseWriter, r *http.Request) error {

	// Render json of our count for the javascript to display
	mu.RLock()
	json := fmt.Sprintf("{\"users\":%d}", len(identifiers))
	mu.RUnlock()
	_, err := w.Write([]byte(json))
	return err
}

// UserCount returns a count of users in the last 5 minutes
func UserCount() int {
	mu.RLock()
	defer mu.RUnlock()
	return len(identifiers)
}

// Clean up users list at intervals
func init() {
	purgeUsers()
}

// purgeUsers clears the users list of users who last acted PurgeInterval ago
func purgeUsers() {

	mu.Lock()
	for k, v := range identifiers {
		purgeTime := time.Now().Add(-PurgeInterval)
		if v.Before(purgeTime) {
			delete(identifiers, k)
		}
	}
	mu.Unlock()

	time.AfterFunc(time.Second*60, purgeUsers)

	//	fmt.Printf("Purged users:%d", UserCount())
}
