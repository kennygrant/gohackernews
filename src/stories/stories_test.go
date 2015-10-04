// Tests for the stories package
package stories

import (
	"testing"
)

// Log a failure message, given msg, expected and result
func fail(t *testing.T, msg string, expected interface{}, result interface{}) {
	t.Fatalf("\n------FAILURE------\nTest failed: %s expected:%v result:%v", msg, expected, result)
}

// Test create of Story
func TestCreateStory(t *testing.T) {

	// Test create, validation, filtering
	// this should insert a story into the db
	// Check the name was assigned properly

	// Check created at was set properly

}

// Test update of Story
func TestUpdateStory(t *testing.T) {

}
