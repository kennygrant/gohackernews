package storyactions

import (
	"fmt"
	"strings"
	"time"

	"github.com/fragmenta/query"
	"github.com/fragmenta/server/schedule"

	"github.com/kennygrant/gohackernews/src/lib/facebook"
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

		// Link to the primary url for this type of story
		url := story.PrimaryURL()

		if strings.HasPrefix(url, "/") {
			url = "https://golangnews.com" + url
		}

		tweet := fmt.Sprintf("%s #golang %s", story.Name, url)

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
	} else {
		context.Logf("#warn no top story found for tweet")
	}

}

// FacebookPostTopStory facebook posts the top story
func FacebookPostTopStory(context schedule.Context) {
	context.Log("#info posting top story facebook")

	// Get the top story
	q := stories.Popular().Limit(1).Order("rank desc, points desc, id desc")

	// Don't fetch old stories
	q.Where("created_at > current_timestamp - interval '6 hours'")

	// Fetch the story
	results, err := stories.FindAll(q)
	if err != nil {
		context.Logf("#error getting top story for fb %s", err)
		return
	}

	if len(results) > 0 {
		story := results[0]
		context.Logf("#info facebook posting %s", story.Name)
		err := facebook.Post(story.Name, story.Url)
		if err != nil {
			context.Logf("#error facebook post top story %s", err)
			return
		}
		// Do not record fb posts - this could lead to duplicates...
		// we should perhaps have a join table for social media posts, rather than dates on stories?
	} else {
		context.Logf("#warn no top story found for fb")
	}
}
