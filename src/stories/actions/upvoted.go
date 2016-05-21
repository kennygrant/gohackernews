package storyactions

import (
	"github.com/fragmenta/query"
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/kennygrant/gohackernews/src/lib/authorise"
	"github.com/kennygrant/gohackernews/src/stories"
)

// HandleUpvoted displays a list of stories the user has upvoted in the past
func HandleUpvoted(context router.Context) error {

	// Build a query
	q := stories.Query().Limit(listLimit)

	// Select only above 0 points,  Order by rank, then points, then name
	q.Where("points > 0").Order("rank desc, points desc, id desc")

	// Select only stories which the user has upvoted
	user := authorise.CurrentUser(context)
	if !user.Anon() {
		// Can we use a join instead?
		v := query.New("votes", "story_id").Select("select story_id as id from votes").Where("user_id=? AND story_id IS NOT NULL AND points > 0", user.Id)

		storyIDs := v.ResultIDs()
		if len(storyIDs) > 0 {
			q.WhereIn("id", storyIDs)
		}
	}

	// Set the offset in pages if we have one
	page := int(context.ParamInt("page"))
	if page > 0 {
		q.Offset(listLimit * page)
	}

	// Fetch the stories
	results, err := stories.FindAll(q)
	if err != nil {
		return router.InternalError(err)
	}

	// Render the template
	view := view.New(context)
	setStoriesMetadata(view, context.Request())
	view.AddKey("page", page)
	view.AddKey("stories", results)
	view.Template("stories/views/index.html.got")

	if context.Param("format") == ".xml" {
		view.Layout("")
		view.Template("stories/views/index.xml.got")
	}

	return view.Render()

}
