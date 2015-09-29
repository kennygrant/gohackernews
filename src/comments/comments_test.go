// Tests for the comments package
package comments

import (
	"testing"
)

// Log a failure message, given msg, expected and result
func fail(t *testing.T, msg string, expected interface{}, result interface{}) {
	t.Fatalf("\n------FAILURE------\nTest failed: %s expected:%v result:%v", msg, expected, result)
}

// Test create of Comment
func TestCreateComment(t *testing.T) {

}

// Test update of Comment
func TestUpdateComment(t *testing.T) {

}
