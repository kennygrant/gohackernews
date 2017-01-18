// Tests for the comments package
package comments

import (
	"testing"

	"github.com/kennygrant/gohackernews/src/lib/resource"
)

var testName = "foo"

func TestSetup(t *testing.T) {
	err := resource.SetupTestDatabase(2)
	if err != nil {
		t.Fatalf("comments: Setup db failed %s", err)
	}

	// Create comments to test with ?

}

// Test Create method
func TestCreateComments(t *testing.T) {
	commentParams := map[string]string{
		"user_name": testName,
		"status":    "100",
	}

	id, err := New().Create(commentParams)
	if err != nil {
		t.Fatalf("comments: Create comment failed :%s", err)
	}

	comment, err := Find(id)
	if err != nil {
		t.Fatalf("comments: Create comment find failed")
	}

	if comment.UserName != testName {
		t.Fatalf("comments: Create comment name failed expected:%s got:%s", testName, comment.UserName)
	}

}

// Test Index (List) method
func TestListComments(t *testing.T) {

	// Get all comments (we should have at least one)
	results, err := FindAll(Query())
	if err != nil {
		t.Fatalf("comments: List no comment found :%s", err)
	}

	if len(results) < 1 {
		t.Fatalf("comments: List no comments found :%s", err)
	}

}

// Test Update method
func TestUpdateComments(t *testing.T) {

	// Get the last comment (created in TestCreateComments above)
	comment, err := FindFirst("user_name=?", testName)
	if err != nil {
		t.Fatalf("comments: Update no comment found :%s", err)
	}

	name := "bar"
	commentParams := map[string]string{"user_name": name}
	err = comment.Update(commentParams)
	if err != nil {
		t.Fatalf("comments: Update comment failed :%s", err)
	}

	// Fetch the comment again from db
	comment, err = Find(comment.ID)
	if err != nil {
		t.Fatalf("comments: Update comment fetch failed :%s", comment.UserName)
	}

	if comment.UserName != name {
		t.Fatalf("comments: Update comment failed :%s", comment.UserName)
	}

}

// TestQuery tests trying to find published resources
func TestQuery(t *testing.T) {

	results, err := FindAll(Published())
	if err != nil {
		t.Fatalf("comments: error getting comments :%s", err)
	}
	if len(results) == 0 {
		t.Fatalf("comments: published comments not found :%s", err)
	}

	results, err = FindAll(Query().Where("id>=? AND id <=?", 0, 100))
	if err != nil || len(results) == 0 {
		t.Fatalf("comments: no comment found :%s", err)
	}
	if len(results) > 2 {
		t.Fatalf("comments: more than one comment found for where :%s", err)
	}

}

// Test Destroy method
func TestDestroyComments(t *testing.T) {

	results, err := FindAll(Query())
	if err != nil || len(results) == 0 {
		t.Fatalf("comments: Destroy no comment found :%s", err)
	}
	comment := results[0]
	count := len(results)

	err = comment.Destroy()
	if err != nil {
		t.Fatalf("comments: Destroy comment failed :%s", err)
	}

	// Check new length of comments returned
	results, err = FindAll(Query())
	if err != nil {
		t.Fatalf("comments: Destroy error getting results :%s", err)
	}

	// length should be one less than previous
	if len(results) != count-1 {
		t.Fatalf("comments: Destroy comment count wrong :%d", len(results))
	}

}

// TestAllowedParams should always return some params
func TestAllowedParams(t *testing.T) {
	if len(AllowedParams()) == 0 {
		t.Fatalf("comments: no allowed params")
	}
}
