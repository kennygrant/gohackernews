// Tests for the users package
package users

import (
	"testing"

	"github.com/fragmenta/query"

	"github.com/kennygrant/gohackernews/src/lib/resource"
)

var testName = "foo"

func TestSetup(t *testing.T) {
	err := resource.SetupTestDatabase(2)
	if err != nil {
		t.Fatalf("users: Setup db failed %s", err)
	}

	// Delete all users first
	_, err = query.ExecSQL("delete from users;")
	if err != nil {
		t.Fatalf("error setting up:%s", err)
	}

	query.ExecSQL("ALTER SEQUENCE users_id_seq RESTART WITH 1;")

}

// Test Create method
func TestCreateUsers(t *testing.T) {
	userParams := map[string]string{
		"name":   testName,
		"status": "100",
	}

	id, err := New().Create(userParams)
	if err != nil {
		t.Fatalf("users: Create user failed :%s", err)
	}

	user, err := Find(id)
	if err != nil {
		t.Fatalf("users: Create user find failed")
	}

	t.Logf("USER FOUND:\n%d-%s-%s", user.ID, user.Name, testName)

	if user.Name != testName {
		t.Fatalf("users: Create user name failed expected:%s got:%s", testName, user.Name)
	}

}

// Test Index (List) method
func TestListUsers(t *testing.T) {

	// Get all users (we should have at least one)
	results, err := FindAll(Query())
	if err != nil {
		t.Fatalf("users: List no user found :%s", err)
	}

	if len(results) < 1 {
		t.Fatalf("users: List no users found :%s", err)
	}

}

// Test Update method
func TestUpdateUsers(t *testing.T) {

	// Get the last user (created in TestCreateUsers above)
	user, err := FindFirst("name=?", testName)
	if err != nil {
		t.Fatalf("users: Update no user found :%s", err)
	}

	name := "bar"
	userParams := map[string]string{"name": name}
	err = user.Update(userParams)
	if err != nil {
		t.Fatalf("users: Update user failed :%s", err)
	}

	// Fetch the user again from db
	user, err = Find(user.ID)
	if err != nil {
		t.Fatalf("users: Update user fetch failed :%s", user.Name)
	}

	if user.Name != name {
		t.Fatalf("users: Update user failed :%s", user.Name)
	}

}

// TestQuery tests trying to find published resources
func TestQuery(t *testing.T) {

	results, err := FindAll(Published())
	if err != nil {
		t.Fatalf("users: error getting users :%s", err)
	}
	if len(results) == 0 {
		t.Fatalf("users: published users not found :%s", err)
	}

	results, err = FindAll(Query().Where("id>=? AND id <=?", 0, 100))
	if err != nil || len(results) == 0 {
		t.Fatalf("users: no user found :%s", err)
	}

}

// Test Destroy method
func TestDestroyUsers(t *testing.T) {

	results, err := FindAll(Query())
	if err != nil || len(results) == 0 {
		t.Fatalf("users: Destroy no user found :%s", err)
	}
	user := results[0]
	count := len(results)

	err = user.Destroy()
	if err != nil {
		t.Fatalf("users: Destroy user failed :%s", err)
	}

	// Check new length of users returned
	results, err = FindAll(Query())
	if err != nil {
		t.Fatalf("users: Destroy error getting results :%s", err)
	}

	// length should be one less than previous
	if len(results) != count-1 {
		t.Fatalf("users: Destroy user count wrong :%d", len(results))
	}

}

// TestAllowedParams should always return some params
func TestAllowedParams(t *testing.T) {
	if len(AllowedParams()) == 0 {
		t.Fatalf("users: no allowed params")
	}
}
