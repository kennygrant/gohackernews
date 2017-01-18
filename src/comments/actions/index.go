package commentactions

import (
	"net/http"

	"github.com/fragmenta/auth/can"
	"github.com/fragmenta/mux"
	"github.com/fragmenta/server"
	"github.com/fragmenta/view"

	"github.com/kennygrant/gohackernews/src/comments"
	"github.com/kennygrant/gohackernews/src/lib/session"
)

// HandleIndex displays a list of comments.
func HandleIndex(w http.ResponseWriter, r *http.Request) error {

	// Authorise list comment
	currentUser := session.CurrentUser(w, r)
	err := can.List(comments.New(), currentUser)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Get the params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Build a query
	q := comments.Query()

	// Order by required order, or default to id asc
	switch params.Get("order") {

	case "1":
		q.Order("created desc")

	case "2":
		q.Order("updated desc")

	default:
		q.Order("id asc")
	}

	// Filter if requested
	filter := params.Get("filter")
	if len(filter) > 0 {
		q.Where("name ILIKE ?", filter)
	}

	// Fetch the comments
	results, err := comments.FindAll(q)
	if err != nil {
		return server.InternalError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.AddKey("filter", filter)
	view.AddKey("comments", results)
	view.AddKey("currentUser", currentUser)
	return view.Render()
}
