package app

import (
	"time"

	"github.com/fragmenta/server/config"
	"github.com/kennygrant/gohackernews/src/lib/twitter"
	"github.com/kennygrant/gohackernews/src/stories/actions"
)

// SetupServices sets up external services from our config file
func SetupServices() {

	// Don't send if not on production server
	if !config.Production() {
		return
	}

	now := time.Now().UTC()

	// Set up twitter if available, and schedule tweets
	if config.Get("twitter_secret") != "" {
		twitter.Setup(config.Get("twitter_key"), config.Get("twitter_secret"), config.Get("twitter_token"), config.Get("twitter_token_secret"))

		tweetTime := time.Date(now.Year(), now.Month(), now.Day(), 6, 0, 0, 0, time.UTC)
		tweetInterval := 5 * time.Hour

		// For testing
		//tweetTime = now.Add(time.Second * 5)

		ScheduleAt(storyactions.TweetTopStory, tweetTime, tweetInterval)
	}
	/*
		// Set up mail
		if config.Get("mail_secret") != "" {
			mail.Setup(config.Get("mail_secret"), config.Get("mail_from"))

			// Schedule emails to go out at 09:00 every day, starting from the next occurance
			emailTime := time.Date(now.Year(), now.Month(), now.Day(), 10, 10, 10, 10, time.UTC)
			emailInterval := 7 * 24 * time.Hour // Send Emails weekly

			// For testing send immediately on launch
			//emailTime = now.Add(time.Second * 2)

			schedule.At(useractions.DailyEmail, context, emailTime, emailInterval)
		}
	*/
}

// ScheduleAt schedules execution for a particular time and at intervals thereafter.
// If interval is 0, the function will be called only once.
// Callers should call close(task) before exiting the app or to stop repeating the action.
func ScheduleAt(f func(), t time.Time, i time.Duration) chan struct{} {
	task := make(chan struct{})
	now := time.Now().UTC()

	// Check that t is not in the past, if it is increment it by interval until it is not
	for now.Sub(t) > 0 {
		t = t.Add(i)
	}

	// We ignore the timer returned by AfterFunc - so no cancelling, perhaps rethink this
	tillTime := t.Sub(now)
	time.AfterFunc(tillTime, func() {
		// Call f at least once at the time specified
		go f()

		// If we have an interval, call it again repeatedly after interval
		// stopping if the caller calls stop(task) on returned channel
		if i > 0 {
			ticker := time.NewTicker(i)
			go func() {
				for {
					select {
					case <-ticker.C:
						go f()
					case <-task:
						ticker.Stop()
						return
					}
				}
			}()
		}
	})

	return task // call close(task) to stop executing the task for repeated tasks
}
