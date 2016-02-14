package useractions

import (
	"time"

	"github.com/fragmenta/query"
	"github.com/fragmenta/server/schedule"

	"github.com/kennygrant/gohackernews/src/lib/mail"
	"github.com/kennygrant/gohackernews/src/stories"
	"github.com/kennygrant/gohackernews/src/users"
)

// DailyEmail sends a daily email to subscribed users with top stories - change this to WeeklyEmail
// before putting it into production
// We should probably only do this for kenny at present
func DailyEmail(context schedule.Context) {
	context.Log("Sending daily email")

	// First fetch our stories over 5 points
	q := stories.Popular()

	// Must be within 7 days
	q.Where("created_at > current_timestamp - interval '7 day'")

	// Order by rank
	q.Order("rank desc, points desc, id desc")

	// Don't fetch stories that have already been mailed
	q.Where("newsletter_at IS NULL")

	// Fetch the stories
	topStories, err := stories.FindAll(q)
	if err != nil {
		context.Logf("#error getting top story tweet %s", err)
		return
	}

	if len(topStories) == 0 {
		context.Logf("#warn no stories found for newsletter")
		return
	}

	// Now fetch our recipient (initially just Kenny as this is in testing)
	recipient, err := users.Find(1)
	if err != nil {
		context.Logf("#error getting email reciipents %s", err)
		return
	}

	var jobStories []*stories.Story

	// Email recipients the stories in question - we should perhaps save in db so that we can
	// have an issue number and always reproduce the digests?
	mailContext := map[string]interface{}{
		"stories": topStories,
		"jobs":    jobStories,
	}
	err = mail.SendOne(recipient.Email, "Go News Digest", "users/views/mail/digest.html.got", mailContext)
	if err != nil {
		context.Logf("#error sending email %s", err)
		return
	}

	// Record that these stories have been mailed in db
	params := map[string]string{"newsletter_at": query.TimeString(time.Now().UTC())}
	err = q.Order("").UpdateAll(params)
	if err != nil {
		context.Logf("#error updating top story newsletter_at %s", err)
		return
	}

}
