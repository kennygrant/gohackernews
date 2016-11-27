// Package twitter providers a wrapper calling the twitter v1.1 API via https
package twitter

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mrjones/oauth"
)

// See https://dev.twitter.com/rest/reference/post/statuses/update

// Twitter config keys - set in the config file or in environment
// these should be set once with setup on app startup, before any request
// assumes caller doesn't want to send with multiple details
var key, secret, accessToken, accessTokenSecret string

// Setup sets our secret keys
func Setup(k, s, at, ats string) error {
	if len(k) == 0 || len(s) == 0 || len(at) == 0 || len(ats) == 0 {
		return fmt.Errorf("#error setting secrets, null value")
	}

	key = k
	secret = s
	accessToken = at
	accessTokenSecret = ats

	return nil
}

// Tweet sends a status update to twitter - returns the response body or error
func Tweet(s string) ([]byte, error) {
	consumer := oauth.NewConsumer(key, secret, oauth.ServiceProvider{})
	//consumer.Debug(true)
	token := &oauth.AccessToken{Token: accessToken, Secret: accessTokenSecret}

	url := "https://api.twitter.com/1.1/statuses/update.json"
	data := map[string]string{"status": s}

	response, err := consumer.Post(url, data, token)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	fmt.Println("Response:", response.StatusCode, response.Status)
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Check for unexpected status codes, and report them
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("#error sending tweet, unexpected status:%d\n\n%s\n", response.StatusCode, body)
	}

	return body, nil
}
