package storyactions

import (
	"fmt"
	"strings"
	"time"

	"github.com/fragmenta/query"
	"github.com/fragmenta/server/config"
	"github.com/fragmenta/server/log"

	"github.com/kennygrant/gohackernews/src/lib/twitter"
	"github.com/kennygrant/gohackernews/src/stories"
)

// TweetTopStory tweets the top story
func TweetTopStory() {
	log.Log(log.Values{"msg": "Sending top story tweet"})

	// Get the top story which has not been tweeted yet, newer than 1 day (we don't look at older stories)
	q := stories.Popular().Limit(1).Order("rank desc, points desc, id desc")

	// Don't fetch old stories
	q.Where("created_at > current_timestamp - interval '3 days'")

	// Don't fetch stories that have already been tweeted
	q.Where("tweeted_at IS NULL")

	// Fetch the stories
	results, err := stories.FindAll(q)
	if err != nil {
		log.Log(log.Values{"message": "stories: error getting top story tweet", "error": err})
		return
	}

	if len(results) > 0 {
		story := results[0]

		TweetStory(story)
	} else {
		log.Log(log.Values{"message": "stories: warning no top story found for tweet"})
	}

}

// TweetStory tweets the given story
func TweetStory(story *stories.Story) {

	// Base url from config
	baseURL := config.Get("root_url")

	// Link to the primary url for this type of story
	url := story.PrimaryURL()

	// Check for relative urls
	if strings.HasPrefix(url, "/") {
		url = baseURL + url
	}

	tweet := fmt.Sprintf("%s #golang %s", story.Name, url)

	// If the tweet will be too long for twitter, use GN url
	if len(tweet) > 140 {
		tweet = fmt.Sprintf("%s #golang %s", story.Name, baseURL+story.ShowURL())
	}

	log.Log(log.Values{"message": "stories: sending tweet", "tweet": tweet})

	_, err := twitter.Tweet(tweet)
	if err != nil {
		log.Log(log.Values{"message": "stories: error tweeting story", "error": err})
		return
	}

	// Record that this story has been tweeted in db
	params := map[string]string{"tweeted_at": query.TimeString(time.Now().UTC())}
	err = story.Update(params)
	if err != nil {
		log.Log(log.Values{"message": "stories: error updating tweeted story", "error": err})
		return
	}

}
