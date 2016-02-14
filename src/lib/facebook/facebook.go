package facebook

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
)

// See https://developers.facebook.com/docs/graph-api/reference/v2.5/page/feed#publish
// See http://nodotcom.org/python-facebook-tutorial.html for access tokens

// Facebook page access token - permissions must already be granted.
var accessToken string

// Setup sets our secret keys
func Setup(secret string) error {
	if len(secret) == 0 {
		return fmt.Errorf("#error setting secrets, null value")
	}
	accessToken = secret
	return nil
}

// Post sends a status update to facebook - returns error
func Post(s string, link string) error {

	url := "https://graph.facebook.com/v2.5/1174411192569505/feed"
	data := map[string]string{"message": s, "link": link, "access_token": accessToken}
	contentType, formData, err := createFormData(data)
	if err != nil {
		return err
	}

	// Now post the form
	req, err := http.NewRequest("POST", url, &formData)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	// Check for unexpected status codes, and report them
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("#error sending facebook status, unexpected status:%d %v\n", resp.StatusCode, resp)
	}

	//fmt.Printf("File sent to: %s %v\n", url, resp.Body)
	fmt.Printf("Post sent to: %s\n", url)

	return nil
}

// createFormData populates a form data object from a map of string-string
func createFormData(data map[string]string) (string, bytes.Buffer, error) {

	// Prepare a new multipart form writer
	var formData bytes.Buffer
	w := multipart.NewWriter(&formData)

	for k, v := range data {
		err := addField(w, k, v)
		if err != nil {
			return "", formData, err
		}
	}

	// Close the writer
	w.Close()

	return w.FormDataContentType(), formData, nil
}

// addField adds a field to the multipart writer
func addField(w *multipart.Writer, k, v string) error {
	fw, err := w.CreateFormField(k)
	if err != nil {
		return err
	}
	_, err = fw.Write([]byte(v))
	if err != nil {
		return err
	}
	return nil
}
