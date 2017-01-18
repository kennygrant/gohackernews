package storyactions

import (
	"net/http"
	"strings"

	"github.com/fragmenta/auth/can"
	"github.com/fragmenta/mux"
	"github.com/fragmenta/server"
	"github.com/fragmenta/server/log"
	"github.com/fragmenta/view"

	"github.com/kennygrant/gohackernews/src/lib/session"
	"github.com/kennygrant/gohackernews/src/stories"
)

// HandleCreateShow serves the create form via GET for stories.
func HandleCreateShow(w http.ResponseWriter, r *http.Request) error {

	story := stories.New()

	// Authorise
	currentUser := session.CurrentUser(w, r)
	err := can.Create(story, currentUser)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.AddKey("story", story)
	view.AddKey("currentUser", currentUser)
	return view.Render()
}

// HandleCreate handles the POST of the create form for stories
func HandleCreate(w http.ResponseWriter, r *http.Request) error {

	story := stories.New()

	// Check the authenticity token
	err := session.CheckAuthenticity(w, r)
	if err != nil {
		return err
	}

	// Get user details
	user := session.CurrentUser(w, r)
	ip := getUserIP(r)

	// Authorise
	err = can.Create(story, user)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Setup context
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Check permissions - if not logged in and above 1 points, redirect to error
	if !user.CanSubmit() {
		return server.NotAuthorizedError(nil, "Sorry", "You need to be registered and have more than 1 points to submit stories.")
	}

	// TODO refactor this and put it in the model instead as a method

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
	// for now only strip on medium urls
	if strings.Contains(url, "#") && strings.Contains(url, "medium.com") {
		url = strings.Split(url, "#")[0]
	}

	// Rewrite mobile youtube links
	if strings.HasPrefix(url, "https://m.youtube.com") {
		url = strings.Replace(url, "https://m.youtube.com", "https://www.youtube.com", 1)
	}

	params.Set("url", []string{url})

	// Check that no story with this url already exists
	q := stories.Where("url=?", url)
	duplicates, err := stories.FindAll(q)
	if err != nil {
		return server.InternalError(err)
	}

	// If we have a duplicate story, with the same non-null url, upvote or reject
	if len(duplicates) > 0 && url != "" {
		story = duplicates[0]

		// Check we have no votes already from this user, if we do fail
		if storyHasUserVote(story, user) {
			return server.NotAuthorizedError(err, "Vote Failed", "You have already submitted this story.")
		}

		// Add a point to dupe and show the story
		addStoryVote(story, user, ip, 1)
		return server.Redirect(w, r, story.ShowURL())
	}

	// Set a few params
	params.SetInt("points", 1)
	params.SetInt("user_id", user.ID)
	params.SetString("user_name", user.Name)

	// Clean params according to role
	accepted := stories.AllowedParams()
	if user.Admin() {
		accepted = stories.AllowedParamsAdmin()
	}

	// Validate the params, removing any we don't accept
	storyParams := story.ValidateParams(params.Map(), accepted)

	ID, err := story.Create(storyParams)
	if err != nil {
		return err // Create returns a router.Error
	}

	// Log creation
	log.Info(log.V{"msg": "Created story", "story_id": ID})

	// Redirect to the new story
	story, err = stories.Find(ID)
	if err != nil {
		return server.InternalError(err)
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

	return server.Redirect(w, r, story.IndexURL())
}
