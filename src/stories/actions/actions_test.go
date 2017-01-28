package storyactions

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/fragmenta/mux"
	"github.com/fragmenta/query"

	"github.com/kennygrant/gohackernews/src/lib/resource"
	"github.com/kennygrant/gohackernews/src/stories"
)

// names is used to test setting and getting the first string field of the story.
var names = []string{"foo", "bar"}

// testSetup performs setup for integration tests
// using the test database, real views, and mock authorisation
// If we can run this once for global tests it might be more efficient?
func TestSetup(t *testing.T) {
	err := resource.SetupTestDatabase(3)
	if err != nil {
		fmt.Printf("stories: Setup db failed %s", err)
	}

	// Set up mock auth
	resource.SetupAuthorisation()

	// Load templates for rendering
	resource.SetupView(3)

	router := mux.New()
	mux.SetDefault(router)

	// FIXME - Need to write routes out here again, but without pkg prefix
	// Any neat way to do this instead? We'd need a separate routes package under app...
	router.Add("/stories", nil)
	router.Add("/stories/create", nil)
	router.Add("/stories/create", nil).Post()
	router.Add("/stories/login", nil)
	router.Add("/stories/login", nil).Post()
	router.Add("/stories/login", nil).Post()
	router.Add("/stories/logout", nil).Post()
	router.Add("/stories/{id:\\d+}/update", nil)
	router.Add("/stories/{id:\\d+}/update", nil).Post()
	router.Add("/stories/{id:\\d+}/destroy", nil).Post()
	router.Add("/stories/{id:\\d+}", nil)

	// Delete all stories to ensure we get consistent results
	query.ExecSQL("delete from stories;")
	query.ExecSQL("ALTER SEQUENCE stories_id_seq RESTART WITH 1;")

	// Delete all users to ensure we get consistent results?
	_, err = query.ExecSQL("delete from users;")
	if err != nil {
		t.Fatalf("error setting up:%s", err)
	}
	// Insert a test admin user for checking logins - never delete as will
	// be required for other resources testing
	_, err = query.ExecSQL("INSERT INTO users (id,email,name,points,status,role,password_hash) VALUES(1,'example@example.com','admin',10,100,100,'$2a$10$2IUzpI/yH0Xc.qs9Z5UUL.3f9bqi0ThvbKs6Q91UOlyCEGY8hdBw6');")
	if err != nil {
		t.Fatalf("error setting up:%s", err)
	}
	// Insert user to delete
	_, err = query.ExecSQL("INSERT INTO users (id,email,name,points,status,role,password_hash) VALUES(2,'example@example.com','test',10,100,0,'$2a$10$2IUzpI/yH0Xc.qs9Z5UUL.3f9bqi0ThvbKs6Q91UOlyCEGY8hdBw6');")
	if err != nil {
		t.Fatalf("error setting up:%s", err)
	}

	query.ExecSQL("ALTER SEQUENCE users_id_seq RESTART WITH 1;")

}

// Test GET /stories/create
func TestShowCreateStories(t *testing.T) {

	// Setup request and recorder
	r := httptest.NewRequest("GET", "/stories/create", nil)
	w := httptest.NewRecorder()

	// Set up story session cookie for admin story above
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("storyactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleCreateShow(w, r)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Fatalf("storyactions: error handling HandleCreateShow %s", err)
	}

	// Test the body for a known pattern
	pattern := "resource-update-form"
	if !strings.Contains(w.Body.String(), pattern) {
		t.Fatalf("storyactions: unexpected response for HandleCreateShow expected:%s got:%s", pattern, w.Body.String())
	}

}

// Test POST /stories/create
func TestCreateStories(t *testing.T) {

	form := url.Values{}
	form.Add("name", names[0])
	body := strings.NewReader(form.Encode())

	r := httptest.NewRequest("POST", "/stories/create", body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Set up story session cookie for admin story
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("storyactions: error setting session %s", err)
	}

	// Run the handler to update the story
	err = HandleCreate(w, r)
	if err != nil {
		t.Fatalf("storyactions: error handling HandleCreate %s", err)
	}

	// Test we get a redirect after update (to the story concerned)
	if w.Code != http.StatusFound {
		t.Fatalf("storyactions: unexpected response code for HandleCreate expected:%d got:%d", http.StatusFound, w.Code)
	}

	// Check the story name is in now value names[1]
	allStories, err := stories.FindAll(stories.Query().Order("id desc"))
	if err != nil || len(allStories) == 0 {
		t.Fatalf("storyactions: error finding created story %s", err)
	}
	newStories := allStories[0]
	if newStories.ID != 1 || newStories.Name != names[0] {
		t.Fatalf("storyactions: error with created story values: %v %s", newStories.ID, newStories.Name)
	}
}

// Test GET /stories
func TestListStories(t *testing.T) {

	// Setup request and recorder
	r := httptest.NewRequest("GET", "/stories", nil)
	w := httptest.NewRecorder()

	// Set up story session cookie for admin story above
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("storyactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleIndex(w, r)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Fatalf("storyactions: error handling HandleIndex %s", err)
	}

	// Test the body for a known pattern
	pattern := `<ul class="sections">`
	if !strings.Contains(w.Body.String(), pattern) {
		t.Fatalf("storyactions: unexpected response for HandleIndex expected:%s got:%s", pattern, w.Body.String())
	}

}

// Test of GET /stories/1
func TestShowStories(t *testing.T) {

	// Setup request and recorder
	r := httptest.NewRequest("GET", "/stories/1", nil)
	w := httptest.NewRecorder()

	// Set up story session cookie for admin story above
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("storyactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleShow(w, r)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Fatalf("storyactions: error handling HandleShow %s", err)
	}

	// Test the body for a known pattern
	pattern := names[0]
	if !strings.Contains(w.Body.String(), names[0]) {
		t.Fatalf("storyactions: unexpected response for HandleShow expected:%s got:%s", pattern, w.Body.String())
	}
}

// Test GET /stories/123/update
func TestShowUpdateStories(t *testing.T) {

	// Setup request and recorder
	r := httptest.NewRequest("GET", "/stories/1/update", nil)
	w := httptest.NewRecorder()

	// Set up story session cookie for admin story above
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("storyactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleUpdateShow(w, r)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Fatalf("storyactions: error handling HandleCreateShow %s", err)
	}

	// Test the body for a known pattern
	pattern := "resource-update-form"
	if !strings.Contains(w.Body.String(), pattern) {
		t.Fatalf("storyactions: unexpected response for HandleCreateShow expected:%s got:%s", pattern, w.Body.String())
	}

}

// Test POST /stories/123/update
func TestUpdateStories(t *testing.T) {

	form := url.Values{}
	form.Add("name", names[1])
	body := strings.NewReader(form.Encode())

	r := httptest.NewRequest("POST", "/stories/1/update", body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Set up story session cookie for admin story
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("storyactions: error setting session %s", err)
	}

	// Run the handler to update the story
	err = HandleUpdate(w, r)
	if err != nil {
		t.Fatalf("storyactions: error handling HandleUpdateStories %s", err)
	}

	// Test we get a redirect after update (to the story concerned)
	if w.Code != http.StatusFound {
		t.Fatalf("storyactions: unexpected response code for HandleUpdateStories expected:%d got:%d", http.StatusFound, w.Code)
	}

	// Check the story name is in now value names[1]
	story, err := stories.Find(1)
	if err != nil {
		t.Fatalf("storyactions: error finding updated story %s", err)
	}
	if story.ID != 1 || story.Name != names[1] {
		t.Fatalf("storyactions: error with updated story values: %v", story)
	}

}

// Test of POST /stories/123/destroy
func TestDeleteStories(t *testing.T) {

	body := strings.NewReader(``)

	// Now test deleting the story created above as admin
	r := httptest.NewRequest("POST", "/stories/1/destroy", body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Set up story session cookie for admin story
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("storyactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleDestroy(w, r)

	// Test the error response is 302 StatusFound
	if err != nil {
		t.Fatalf("storyactions: error handling HandleDestroy %s", err)
	}

	// Test we get a redirect after delete
	if w.Code != http.StatusFound {
		t.Fatalf("storyactions: unexpected response code for HandleDestroy expected:%d got:%d", http.StatusFound, w.Code)
	}
	// Now test as anon
	r = httptest.NewRequest("POST", "/stories/1/destroy", body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()

	// Run the handler to test failure as anon
	err = HandleDestroy(w, r)
	if err == nil { // failure expected
		t.Fatalf("storyactions: unexpected response for HandleDestroy as anon, expected failure")
	}

}
