package storyactions

import (
	"fmt"
	"strings"

	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	"github.com/kennygrant/gohackernews/src/lib/authorise"
	"github.com/kennygrant/gohackernews/src/stories"
)

// HandleCreateShow serves the create form via GET for stories
func HandleCreateShow(context router.Context) error {

	// Authorise
	err := authorise.Path(context)
	if err != nil {
		// When not allowed to post stories, redirect to register screen
		router.Redirect(context, "/users/create")
	}

	// Render the template
	view := view.New(context)
	story := stories.New()
	view.AddKey("story", story)
	view.AddKey("meta_title", "Go Hacker News Submit")

	return view.Render()
}

// HandleCreate handles the POST of the create form for stories
func HandleCreate(context router.Context) error {

	// Check csrf token
	err := authorise.AuthenticityToken(context)
	if err != nil {
		return router.NotAuthorizedError(err)
	}

	// Check permissions - if not logged in and above 1 points, redirect to error
	if !authorise.CurrentUser(context).CanSubmit() {
		return router.NotAuthorizedError(nil, "Sorry", "You need to be registered and have more than 1 points to submit stories.")
	}

	// Get params
	params, err := context.Params()
	if err != nil {
		return router.InternalError(err)
	}

	// Get user details
	user := authorise.CurrentUser(context)
	ip := getUserIP(context)

	// Process urls
	url := params.Get("url")

	// Strip trailing slashes on url before comparisons
	if strings.HasSuffix(url, "/") {
		url = strings.Trim(url, "/")
	}

	// Strip ?utm_source etc - remove all after ?utm_source
	if strings.Contains(url, "?utm_") {
		url = strings.Split(url, "?utm_")[0]
	}

	// Strip url fragments (For example trailing # on medium urls)
	if strings.Contains(url, "#") {
		url = strings.Split(url, "#")[0]
	}

	// Rewrite mobile youtube links
	if strings.HasPrefix(url, "https://m.youtube.com") {
		url = strings.Replace(url, "https://m.youtube.com", "https://www.youtube.com", 1)
	}

	params.Set("url", url)

	// Check that no story with this url already exists
	q := stories.Where("url=?", url)
	duplicates, err := stories.FindAll(q)
	if err != nil {
		return router.InternalError(err)
	}

	if len(duplicates) > 0 {
		story := duplicates[0]

		// Check we have no votes already from this user, if we do fail
		if storyHasUserVote(story, user) {
			return router.NotAuthorizedError(err, "Vote Failed", "Sorry you are not allowed to vote twice, nice try Luis!")
		}

		// Add a point to dupe and return
		addStoryVote(story, user, ip, 1)
		return router.Redirect(context, story.URLShow())
	}
	// Clean params according to role
	accepted := stories.AllowedParams()
	if authorise.CurrentUser(context).Admin() {
		accepted = stories.AllowedParamsAdmin()
	}
	cleanedParams := params.Clean(accepted)

	// Set a few params
	cleanedParams["points"] = "1"
	cleanedParams["user_id"] = fmt.Sprintf("%d", user.Id)
	cleanedParams["user_name"] = user.Name

	id, err := stories.Create(cleanedParams)
	if err != nil {
		return err // Create returns a router.Error
	}

	// Log creation
	context.Logf("#info Created story id,%d", id)

	// Redirect to the new story
	story, err := stories.Find(id)
	if err != nil {
		return router.InternalError(err)
	}

	// We need to add a vote to the story here too by adding a join to the new id
	err = recordStoryVote(story, user, ip, +1)
	if err != nil {
		return err
	}

	// Re-rank stories
	err = updateStoriesRank()
	if err != nil {
		return err
	}

	return router.Redirect(context, story.URLIndex())
}
