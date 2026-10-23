package zentests

import (
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

// These tests check that each assertion reports a wrong response. Each one runs the
// assertion against a deliberately wrong response and checks the number of failures.

// TestStatus_WrongCode_Fails checks that a status mismatch is reported.
func TestStatus_WrongCode_Fails(t *testing.T) {
	errs := failuresOf(t, sendString("ok"), func(r *Response) { r.Status(404) })
	assert.Len(t, errs, 1)
}

// TestHasHeader_WrongValue_Fails checks that a header value mismatch is reported.
func TestHasHeader_WrongValue_Fails(t *testing.T) {
	errs := failuresOf(t, sendString("ok"), func(r *Response) { r.HasHeader("Content-Type", "application/json") })
	assert.Len(t, errs, 1)
}

// TestHeaderContains_MissingSubstring_Fails checks that a missing header substring is reported.
func TestHeaderContains_MissingSubstring_Fails(t *testing.T) {
	errs := failuresOf(t, sendString("ok"), func(r *Response) { r.HeaderContains("Content-Type", "image") })
	assert.Len(t, errs, 1)
}

// TestHeaderPresent_Missing_Fails checks that an absent header is reported.
func TestHeaderPresent_Missing_Fails(t *testing.T) {
	errs := failuresOf(t, sendString("ok"), func(r *Response) { r.HeaderPresent("X-Missing") })
	assert.Len(t, errs, 1)
}

// TestHeaderNotPresent_Present_Fails checks that a header that should be absent is reported.
func TestHeaderNotPresent_Present_Fails(t *testing.T) {
	errs := failuresOf(t, sendString("ok"), func(r *Response) { r.HeaderNotPresent("Content-Type") })
	assert.Len(t, errs, 1)
}

// TestHeaderHasValues_Mismatch_Fails checks that a multi-value header mismatch is reported.
func TestHeaderHasValues_Mismatch_Fails(t *testing.T) {
	errs := failuresOf(t, sendString("ok"), func(r *Response) { r.HeaderHasValues("Content-Type", []string{"a", "b"}) })
	assert.NotEmpty(t, errs)
}

// TestIsJSON_PlainText_Fails checks that IsJSON rejects a text response.
func TestIsJSON_PlainText_Fails(t *testing.T) {
	errs := failuresOf(t, sendString("ok"), func(r *Response) { r.IsJSON() })
	assert.Len(t, errs, 1)
}

// TestIsImage_PlainText_Fails checks that IsImage rejects a text response.
func TestIsImage_PlainText_Fails(t *testing.T) {
	errs := failuresOf(t, sendString("ok"), func(r *Response) { r.IsImage() })
	assert.Len(t, errs, 1)
}

// TestCookieHasValues_WrongValue_Fails checks that a cookie value mismatch is reported.
func TestCookieHasValues_WrongValue_Fails(t *testing.T) {
	handler := func(c fiber.Ctx) error {
		c.Cookie(&fiber.Cookie{Name: "session", Value: "abc"})
		return c.SendString("ok")
	}
	errs := failuresOf(t, handler, func(r *Response) {
		r.CookieHasValues(map[string]string{"session": "xyz"})
	})
	assert.Len(t, errs, 1)
}

// TestContains_MissingSubstring_Fails checks that a missing body substring is reported.
func TestContains_MissingSubstring_Fails(t *testing.T) {
	errs := failuresOf(t, sendString("hello"), func(r *Response) { r.Contains("bye") })
	assert.Len(t, errs, 1)
}

// TestNotContains_PresentSubstring_Fails checks that a forbidden body substring is reported.
func TestNotContains_PresentSubstring_Fails(t *testing.T) {
	errs := failuresOf(t, sendString("hello"), func(r *Response) { r.NotContains("hello") })
	assert.Len(t, errs, 1)
}

// TestBodyMatches_NoMatch_Fails checks that a non-matching pattern is reported.
func TestBodyMatches_NoMatch_Fails(t *testing.T) {
	errs := failuresOf(t, sendString("hello"), func(r *Response) { r.BodyMatches(`\d+`) })
	assert.Len(t, errs, 1)
}

// TestEquals_Mismatch_Fails checks that an exact body mismatch is reported.
func TestEquals_Mismatch_Fails(t *testing.T) {
	errs := failuresOf(t, sendString("hello"), func(r *Response) { r.Equals("bye") })
	assert.Len(t, errs, 1)
}

// TestIsEmpty_NonEmptyBody_Fails checks that a non-empty body fails IsEmpty.
func TestIsEmpty_NonEmptyBody_Fails(t *testing.T) {
	errs := failuresOf(t, sendString("hello"), func(r *Response) { r.IsEmpty() })
	assert.Len(t, errs, 1)
}

// TestHasKey_MissingKey_ReportsOnce checks that a missing JSON key is reported once.
func TestHasKey_MissingKey_ReportsOnce(t *testing.T) {
	errs := failuresOf(t, sendRawJSON(`{"a":1}`), func(r *Response) { r.HasKey("b") })
	assert.Len(t, errs, 1)
}

// TestHasKey_ArrayIndexOutOfBounds_Fails checks that an out-of-range array index is reported.
func TestHasKey_ArrayIndexOutOfBounds_Fails(t *testing.T) {
	errs := failuresOf(t, sendRawJSON(`{"items":[1]}`), func(r *Response) { r.HasKey("items.5") })
	assert.Len(t, errs, 1)
}

// TestHasKey_NonNumericIndex_Fails checks that a non-numeric index into an array is reported.
func TestHasKey_NonNumericIndex_Fails(t *testing.T) {
	errs := failuresOf(t, sendRawJSON(`{"items":[1]}`), func(r *Response) { r.HasKey("items.first") })
	assert.Len(t, errs, 1)
}

// TestHasKey_TraverseScalar_Fails checks that a path cannot continue past a scalar value.
func TestHasKey_TraverseScalar_Fails(t *testing.T) {
	errs := failuresOf(t, sendRawJSON(`{"a":1}`), func(r *Response) { r.HasKey("a.b") })
	assert.Len(t, errs, 1)
}

// TestHas_StringMismatch_Fails checks that a string value mismatch is reported.
func TestHas_StringMismatch_Fails(t *testing.T) {
	errs := failuresOf(t, sendRawJSON(`{"name":"John"}`), func(r *Response) { r.Has("name", "Jane") })
	assert.Len(t, errs, 1)
}

// TestHasString_WrongValue_Fails checks that a string mismatch is reported.
func TestHasString_WrongValue_Fails(t *testing.T) {
	errs := failuresOf(t, sendRawJSON(`{"name":"John"}`), func(r *Response) { r.HasString("name", "Jane") })
	assert.Len(t, errs, 1)
}

// TestHasString_NumberValue_Fails checks that a number does not match an expected string.
func TestHasString_NumberValue_Fails(t *testing.T) {
	errs := failuresOf(t, sendRawJSON(`{"n":1}`), func(r *Response) { r.HasString("n", "1") })
	assert.Len(t, errs, 1)
}

// TestHasBool_WrongValue_Fails checks that a boolean mismatch is reported.
func TestHasBool_WrongValue_Fails(t *testing.T) {
	errs := failuresOf(t, sendRawJSON(`{"ok":false}`), func(r *Response) { r.HasBool("ok", true) })
	assert.Len(t, errs, 1)
}

// TestHasInt_StringValue_Fails checks that a non-numeric value is reported as a type mismatch.
func TestHasInt_StringValue_Fails(t *testing.T) {
	errs := failuresOf(t, sendRawJSON(`{"n":"42"}`), func(r *Response) { r.HasInt("n", 42) })
	assert.Len(t, errs, 1)
}

// TestHasInt_WrongValue_Fails checks that a different integer is reported.
func TestHasInt_WrongValue_Fails(t *testing.T) {
	errs := failuresOf(t, sendRawJSON(`{"n":41}`), func(r *Response) { r.HasInt("n", 42) })
	assert.Len(t, errs, 1)
}

// TestHasFloat_StringValue_Fails checks that a non-numeric value is reported as a type mismatch.
func TestHasFloat_StringValue_Fails(t *testing.T) {
	errs := failuresOf(t, sendRawJSON(`{"score":"9"}`), func(r *Response) { r.HasFloat("score", 9.0) })
	assert.Len(t, errs, 1)
}

// TestHasFloat_WrongValue_Fails checks that a different float is reported.
func TestHasFloat_WrongValue_Fails(t *testing.T) {
	errs := failuresOf(t, sendRawJSON(`{"score":9.5}`), func(r *Response) { r.HasFloat("score", 9.0) })
	assert.Len(t, errs, 1)
}

// TestMatchesRegex_NoMatch_Fails checks that a non-matching value is reported.
func TestMatchesRegex_NoMatch_Fails(t *testing.T) {
	errs := failuresOf(t, sendRawJSON(`{"email":"nope"}`), func(r *Response) {
		r.MatchesRegex("email", `^[\w.-]+@[\w.-]+$`)
	})
	assert.Len(t, errs, 1)
}

// TestJSONHasValues_Mismatch_Fails checks that a mismatch on any listed key is reported.
func TestJSONHasValues_Mismatch_Fails(t *testing.T) {
	errs := failuresOf(t, sendRawJSON(`{"name":"John","age":30}`), func(r *Response) {
		r.JSONHasValues(map[string]any{"name": "John", "age": 31})
	})
	assert.Len(t, errs, 1)
}

// TestIsNull_NonNullValue_Fails checks that a non-null value fails IsNull.
func TestIsNull_NonNullValue_Fails(t *testing.T) {
	errs := failuresOf(t, sendRawJSON(`{"a":1}`), func(r *Response) { r.IsNull("a") })
	assert.Len(t, errs, 1)
}

// TestIsNotNull_NullValue_Fails checks that a null value fails IsNotNull.
func TestIsNotNull_NullValue_Fails(t *testing.T) {
	errs := failuresOf(t, sendRawJSON(`{"a":null}`), func(r *Response) { r.IsNotNull("a") })
	assert.Len(t, errs, 1)
}

// TestArrayLength_WrongLength_Fails checks that an array length mismatch is reported.
func TestArrayLength_WrongLength_Fails(t *testing.T) {
	errs := failuresOf(t, sendRawJSON(`{"items":[1,2]}`), func(r *Response) { r.ArrayLength("items", 3) })
	assert.Len(t, errs, 1)
}

// TestJSON_InvalidJSON_Fails checks that a body that is not JSON is reported.
func TestJSON_InvalidJSON_Fails(t *testing.T) {
	errs := failuresOf(t, sendRawJSON(`{not json`), func(r *Response) { r.JSON() })
	assert.Len(t, errs, 1)
}
