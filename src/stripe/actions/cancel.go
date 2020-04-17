package stripeactions

import (
	"net/http"

	"github.com/fragmenta/mux/log"
	"github.com/fragmenta/view"

	"github.com/kennygrant/gohackernews/src/lib/session"
)

// HandleShowPayCancel displays a page for canceleld payments
func HandleShowPayCancel(w http.ResponseWriter, r *http.Request) error {

	// Find logged in user (if any)
	currentUser := session.CurrentUser(w, r)

	log.Printf("pay: page thanks shown for user:%d", currentUser.ID)

	// Render the template
	view := view.NewRenderer(w, r)
	view.AddKey("currentUser", currentUser)
	return view.Render()
}
