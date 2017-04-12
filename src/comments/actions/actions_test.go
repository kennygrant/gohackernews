package commentactions

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/fragmenta/mux"
	"github.com/fragmenta/query"

	"github.com/kennygrant/gohackernews/src/comments"
	"github.com/kennygrant/gohackernews/src/lib/resource"
)

// names is used to test setting and getting the first string field of the comment.
var names = []string{"foo", "bar"}

// testSetup performs setup for integration tests
// using the test database, real views, and mock authorisation
// If we can run this once for global tests it might be more efficient?
func TestSetup(t *testing.T) {
	err := resource.SetupTestDatabase(3)
	if err != nil {
		fmt.Printf("comments: Setup db failed %s", err)
	}

	// Set up mock auth
	resource.SetupAuthorisation()

	// Load templates for rendering
	resource.SetupView(3)

	router := mux.New()
	mux.SetDefault(router)

	// FIXME - Need to write routes out here again, but without pkg prefix
	// Any neat way to do this instead? We'd need a separate routes package under app...
	router.Add("/comments", nil)
	router.Add("/comments/create", nil)
	router.Add("/comments/create", nil).Post()
	router.Add("/comments/login", nil)
	router.Add("/comments/login", nil).Post()
	router.Add("/comments/login", nil).Post()
	router.Add("/comments/logout", nil).Post()
	router.Add("/comments/{id:\\d+}/update", nil)
	router.Add("/comments/{id:\\d+}/update", nil).Post()
	router.Add("/comments/{id:\\d+}/destroy", nil).Post()
	router.Add("/comments/{id:\\d+}", nil)

	// Delete all comments to ensure we get consistent results
	query.ExecSQL("delete from comments;")
	query.ExecSQL("ALTER SEQUENCE comments_id_seq RESTART WITH 1;")

	// Delete all users to ensure we get consistent results?
	_, err = query.ExecSQL("delete from users;")
	if err != nil {
		t.Fatalf("error setting up:%s", err)
	}
	// Insert a test admin user for checking logins - never delete as will
	// be required for other resources testing
	_, err = query.ExecSQL("INSERT INTO users (id,email,name,points,status,role,password_hash) VALUES(1,'example@example.com','admin',100,100,100,'$2a$10$2IUzpI/yH0Xc.qs9Z5UUL.3f9bqi0ThvbKs6Q91UOlyCEGY8hdBw6');")
	if err != nil {
		t.Fatalf("error setting up:%s", err)
	}
	// Insert user to delete
	_, err = query.ExecSQL("INSERT INTO users (id,email,name,points,status,role,password_hash) VALUES(2,'example2@example.com','test',100,100,0,'$2a$10$2IUzpI/yH0Xc.qs9Z5UUL.3f9bqi0ThvbKs6Q91UOlyCEGY8hdBw6');")
	if err != nil {
		t.Fatalf("error setting up:%s", err)
	}

	query.ExecSQL("ALTER SEQUENCE users_id_seq RESTART WITH 1;")

	// Insert story as a parent
	_, err = query.ExecSQL("INSERT INTO stories (id,name,points) VALUES(1,'Story',100);")
	if err != nil {
		t.Fatalf("error setting up story:%s", err)
	}

}

// Test GET /comments/create
func TestShowCreateComments(t *testing.T) {

	// Setup request and recorder
	r := httptest.NewRequest("GET", "/comments/create", nil)
	w := httptest.NewRecorder()

	// Set up comment session cookie for admin comment above
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("commentactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleCreateShow(w, r)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Fatalf("commentactions: error handling HandleCreateShow %s", err)
	}

	// Test the body for a known pattern
	pattern := "resource-update-form"
	if !strings.Contains(w.Body.String(), pattern) {
		t.Fatalf("commentactions: unexpected response for HandleCreateShow expected:%s got:%s", pattern, w.Body.String())
	}

}

// Test POST /comments/create
func TestCreateComments(t *testing.T) {

	form := url.Values{}
	form.Add("user_name", names[0])
	form.Add("user_id", "1")
	form.Add("story_id", "1")
	form.Add("text", "foo bar comment")
	body := strings.NewReader(form.Encode())

	r := httptest.NewRequest("POST", "/comments/create", body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Set up comment session cookie for admin comment
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("commentactions: error setting session %s", err)
	}

	// Run the handler to update the comment
	err = HandleCreate(w, r)
	if err != nil {
		t.Fatalf("commentactions: error handling HandleCreate %s", err)
	}

	// Test we get a redirect after update (to the comment concerned)
	if w.Code != http.StatusFound {
		t.Fatalf("commentactions: unexpected response code for HandleCreate expected:%d got:%d", http.StatusFound, w.Code)
	}

	// Check the comment name is in now value names[1]
	allComments, err := comments.FindAll(comments.Query().Order("id desc"))
	if err != nil || len(allComments) == 0 {
		t.Fatalf("commentactions: error finding created comment %s", err)
	}
	newComments := allComments[0]
	if newComments.ID != 1 || newComments.UserName != "admin" { // user name of admin user used to create
		t.Fatalf("commentactions: error with created comment values: %v %s", newComments.ID, newComments.UserName)
	}
}

// Test GET /comments
func TestListComments(t *testing.T) {

	// Setup request and recorder
	r := httptest.NewRequest("GET", "/comments", nil)
	w := httptest.NewRecorder()

	// Set up comment session cookie for admin comment above
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("commentactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleIndex(w, r)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Fatalf("commentactions: error handling HandleIndex %s", err)
	}

	// Test the body for a known pattern
	pattern := `<ul class="comments">`
	if !strings.Contains(w.Body.String(), pattern) {
		t.Fatalf("commentactions: unexpected response for HandleIndex expected:%s got:%s", pattern, w.Body.String())
	}

}

// Test of GET /comments/1
func TestShowComments(t *testing.T) {

	// Setup request and recorder
	r := httptest.NewRequest("GET", "/comments/1", nil)
	w := httptest.NewRecorder()

	// Set up comment session cookie for admin comment above
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("commentactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleShow(w, r)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Fatalf("commentactions: error handling HandleShow %s", err)
	}

	// Test the body for a known pattern
	pattern := names[0]
	if !strings.Contains(w.Body.String(), names[0]) {
		t.Fatalf("commentactions: unexpected response for HandleShow expected:%s got:%s", pattern, w.Body.String())
	}
}

// Test GET /comments/123/update
func TestShowUpdateComments(t *testing.T) {

	// Setup request and recorder
	r := httptest.NewRequest("GET", "/comments/1/update", nil)
	w := httptest.NewRecorder()

	// Set up comment session cookie for admin comment above
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("commentactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleUpdateShow(w, r)

	// Test the error response
	if err != nil || w.Code != http.StatusOK {
		t.Fatalf("commentactions: error handling HandleCreateShow %s", err)
	}

	// Test the body for a known pattern
	pattern := "resource-update-form"
	if !strings.Contains(w.Body.String(), pattern) {
		t.Fatalf("commentactions: unexpected response for HandleCreateShow expected:%s got:%s", pattern, w.Body.String())
	}

}

// Test POST /comments/123/update
func TestUpdateComments(t *testing.T) {

	form := url.Values{}
	form.Add("user_name", names[1])
	body := strings.NewReader(form.Encode())

	r := httptest.NewRequest("POST", "/comments/1/update", body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Set up comment session cookie for admin comment
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("commentactions: error setting session %s", err)
	}

	// Run the handler to update the comment
	err = HandleUpdate(w, r)
	if err != nil {
		t.Fatalf("commentactions: error handling HandleUpdateComments %s", err)
	}

	// Test we get a redirect after update (to the comment concerned)
	if w.Code != http.StatusFound {
		t.Fatalf("commentactions: unexpected response code for HandleUpdateComments expected:%d got:%d", http.StatusFound, w.Code)
	}

	// Check the comment name is in now value names[1]
	comment, err := comments.Find(1)
	if err != nil {
		t.Fatalf("commentactions: error finding updated comment %s", err)
	}
	if comment.ID != 1 || comment.UserName != names[1] {
		t.Fatalf("commentactions: error with updated comment values: %v", comment)
	}

}

// Test of POST /comments/123/destroy
func TestDeleteComments(t *testing.T) {

	body := strings.NewReader(``)

	// Now test deleting the comment created above as admin
	r := httptest.NewRequest("POST", "/comments/1/destroy", body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// Set up comment session cookie for admin comment
	err := resource.AddUserSessionCookie(w, r, 1)
	if err != nil {
		t.Fatalf("commentactions: error setting session %s", err)
	}

	// Run the handler
	err = HandleDestroy(w, r)

	// Test the error response is 302 StatusFound
	if err != nil {
		t.Fatalf("commentactions: error handling HandleDestroy %s", err)
	}

	// Test we get a redirect after delete
	if w.Code != http.StatusFound {
		t.Fatalf("commentactions: unexpected response code for HandleDestroy expected:%d got:%d", http.StatusFound, w.Code)
	}
	// Now test as anon
	r = httptest.NewRequest("POST", "/comments/1/destroy", body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w = httptest.NewRecorder()

	// Run the handler to test failure as anon
	err = HandleDestroy(w, r)
	if err == nil { // failure expected
		t.Fatalf("commentactions: unexpected response for HandleDestroy as anon, expected failure")
	}

}
