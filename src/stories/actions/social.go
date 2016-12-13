package storyactions

import (
	"fmt"
	"strings"
	"time"

	"github.com/fragmenta/query"
	"github.com/fragmenta/server/schedule"

	"github.com/kennygrant/gohackernews/src/lib/twitter"
	"github.com/kennygrant/gohackernews/src/stories"
)

// TweetTopStory tweets the top story
func TweetTopStory(context schedule.Context) {
	context.Log("Sending top story tweet")

	// Get the top story which has not been tweeted yet, newer than 1 day (we don't look at older stories)
	q := stories.Popular().Limit(1).Order("rank desc, points desc, id desc")

	// Don't fetch old stories - at some point soon this can come down to 1 day
	// as all older stories will have been tweeted
	// For now omit this as we have a backlog of old unposted stories
	// q.Where("created_at > current_timestamp - interval '60 days'")

	// Don't fetch stories that have already been tweeted
	q.Where("tweeted_at IS NULL")

	// Fetch the stories
	results, err := stories.FindAll(q)
	if err != nil {
		context.Logf("#error getting top story tweet %s", err)
		return
	}

	if len(results) > 0 {
		story := results[0]

		TweetStory(context, story)
	} else {
		context.Logf("#warn no top story found for tweet")
	}

}

// TweetStory tweets the given story
func TweetStory(context schedule.Context, story *stories.Story) {

	// Base url from config
	baseURL := context.Config("root_url")

	// Link to the primary url for this type of story
	url := story.PrimaryURL()

	// Check for relative urls
	if strings.HasPrefix(url, "/") {
		url = baseURL + url
	}

	tweet := fmt.Sprintf("%s #golang %s", story.Name, url)

	// If the tweet will be too long for twitter, use GN url
	if len(tweet) > 140 {
		tweet = fmt.Sprintf("%s #golang %s", story.Name, baseURL+story.URLShow())
	}

	context.Logf("#info sending tweet:%s", tweet)

	_, err := twitter.Tweet(tweet)
	if err != nil {
		context.Logf("#error tweeting top story %s", err)
		return
	}

	// Record that this story has been tweeted in db
	params := map[string]string{"tweeted_at": query.TimeString(time.Now().UTC())}
	err = story.Update(params)
	if err != nil {
		context.Logf("#error updating top story tweet %s", err)
		return
	}

}
