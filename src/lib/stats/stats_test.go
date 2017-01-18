package stats

import (
	"net/http/httptest"
	"testing"
)

// TestStats tests our options are functional when embedded in a resource.
func TestStats(t *testing.T) {

	r := httptest.NewRequest("GET", "/", nil)
	c := UserCount()
	RegisterHit(r)
	newc := UserCount()
	if newc <= c {
		t.Errorf("Stats count incorrect")
	}

	purgeUsers()

	if newc != UserCount() {
		t.Errorf("Stats count incorrect after purge")
	}

	w := httptest.NewRecorder()
	HandleUserCount(w, r)

	// Test recorded user count

}
