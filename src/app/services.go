package app

import (
	"github.com/fragmenta/server/config"
	//	"github.com/fragmenta/server/schedule"
)

// SetupServices sets up external services from our config file
func SetupServices() {

	// Don't send if not on production server
	if !config.Production() {
		return
	}
	/*
		context := schedule.NewContext(server.Logger, server)

		now := time.Now().UTC()

		// Set up twitter if available, and schedule tweets
		if config["twitter_secret"] != "" {
			twitter.Setup(config["twitter_key"], config["twitter_secret"], config["twitter_token"], config["twitter_token_secret"])

			tweetTime := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, time.UTC)
			tweetInterval := 5 * time.Hour

			// For testing
			//tweetTime = now.Add(time.Second * 5)

			schedule.At(storyactions.TweetTopStory, context, tweetTime, tweetInterval)
		}

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
