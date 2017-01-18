package storyactions

import (
	"net/http"
	"strings"

	"github.com/fragmenta/mux"
	"github.com/fragmenta/query"
	"github.com/fragmenta/server"
	"github.com/fragmenta/server/config"
	"github.com/fragmenta/view"

	"github.com/kennygrant/gohackernews/src/lib/session"
	"github.com/kennygrant/gohackernews/src/lib/stats"
	"github.com/kennygrant/gohackernews/src/stories"
)

// HandleListUpvoted displays a list of stories the user has upvoted in the past
func HandleListUpvoted(w http.ResponseWriter, r *http.Request) error {
	stats.RegisterHit(r)

	// Build a query
	q := stories.Query().Limit(listLimit)

	// Select only above 0 points,  Order by rank, then points, then name
	q.Where("points > 0").Order("rank desc, points desc, id desc")

	// Select only stories which the user has upvoted
	user := session.CurrentUser(w, r)
	if !user.Anon() {
		// Can we use a join instead?
		v := query.New("votes", "story_id").Select("select story_id as id from votes").Where("user_id=? AND story_id IS NOT NULL AND points > 0", user.ID)

		storyIDs := v.ResultIDs()
		if len(storyIDs) > 0 {
			q.WhereIn("id", storyIDs)
		}
	}

	// Fetch the  params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Set the offset in pages if we have one
	page := int(params.GetInt("page"))
	if page > 0 {
		q.Offset(listLimit * page)
	}

	// Fetch the stories
	results, err := stories.FindAll(q)
	if err != nil {
		return server.InternalError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)

	view.AddKey("page", page)
	view.AddKey("stories", results)
	view.AddKey("pubdate", storiesModTime(results))
	view.AddKey("meta_title", "Stories you have upvoted")
	view.AddKey("meta_desc", config.Get("meta_desc"))
	view.AddKey("meta_keywords", config.Get("meta_keywords"))
	view.AddKey("meta_rss", storiesXMLPath(w, r))
	view.Template("stories/views/index.html.got")
	view.AddKey("currentUser", session.CurrentUser(w, r))

	if strings.HasSuffix(r.URL.Path, ".xml") {
		view.Layout("")
		view.Template("stories/views/index.xml.got")
	}

	return view.Render()

}
