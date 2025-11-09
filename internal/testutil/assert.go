package testutil

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"testing"
)

func AssertEqual(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		if len(msgAndArgs) > 0 {
			t.Errorf("Expected %v, got %v. %v", expected, actual, msgAndArgs[0])
		} else {
			t.Errorf("Expected %v, got %v", expected, actual)
		}
	}
}

func AssertNotEqual(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if reflect.DeepEqual(expected, actual) {
		if len(msgAndArgs) > 0 {
			t.Errorf("Expected not to be %v. %v", expected, msgAndArgs[0])
		} else {
			t.Errorf("Expected not to be %v", expected)
		}
	}
}

func AssertNil(t *testing.T, value interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if value != nil && !reflect.ValueOf(value).IsNil() {
		if len(msgAndArgs) > 0 {
			t.Errorf("Expected nil, got %v. %v", value, msgAndArgs[0])
		} else {
			t.Errorf("Expected nil, got %v", value)
		}
	}
}

func AssertNotNil(t *testing.T, value interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	if value == nil || (reflect.ValueOf(value).Kind() == reflect.Ptr && reflect.ValueOf(value).IsNil()) {
		if len(msgAndArgs) > 0 {
			t.Errorf("Expected not nil. %v", msgAndArgs[0])
		} else {
			t.Errorf("Expected not nil, got nil")
		}
	}
}

func AssertStatusCode(t *testing.T, expected int, resp *http.Response) {
	t.Helper()
	if resp.StatusCode != expected {
		t.Errorf("Expected status code %d, got %d", expected, resp.StatusCode)
	}
}

func AssertContains(t *testing.T, haystack, needle string, msgAndArgs ...interface{}) {
	t.Helper()
	if len(haystack) == 0 || len(needle) == 0 {
		if len(msgAndArgs) > 0 {
			t.Errorf("String not found. %v", msgAndArgs[0])
		} else {
			t.Errorf("Expected '%s' to contain '%s'", haystack, needle)
		}
		return
	}
	found := false
	for i := 0; i <= len(haystack)-len(needle); i++ {
		if haystack[i:i+len(needle)] == needle {
			found = true
			break
		}
	}
	if !found {
		if len(msgAndArgs) > 0 {
			t.Errorf("Expected '%s' to contain '%s'. %v", haystack, needle, msgAndArgs[0])
		} else {
			t.Errorf("Expected '%s' to contain '%s'", haystack, needle)
		}
	}
}

func AssertEmpty(t *testing.T, value interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	isEmpty := false

	if value == nil {
		isEmpty = true
	} else {
		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.Array, reflect.Slice, reflect.Map, reflect.String:
			isEmpty = v.Len() == 0
		}
	}

	if !isEmpty {
		if len(msgAndArgs) > 0 {
			t.Errorf("Expected empty, got %v. %v", value, msgAndArgs[0])
		} else {
			t.Errorf("Expected empty, got %v", value)
		}
	}
}

func AssertNotEmpty(t *testing.T, value interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	isEmpty := true

	if value == nil {
		isEmpty = true
	} else {
		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.Array, reflect.Slice, reflect.Map, reflect.String:
			isEmpty = v.Len() == 0
		default:
			isEmpty = false
		}
	}

	if isEmpty {
		if len(msgAndArgs) > 0 {
			t.Errorf("Expected not empty. %v", msgAndArgs[0])
		} else {
			t.Errorf("Expected not empty, got empty")
		}
	}
}

func ParseJSONResponse(t *testing.T, body io.Reader, target interface{}) {
	t.Helper()
	if err := json.NewDecoder(body).Decode(target); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}
}

func AssertLen(t *testing.T, value interface{}, expectedLen int, msgAndArgs ...interface{}) {
	t.Helper()
	v := reflect.ValueOf(value)
	actualLen := 0

	switch v.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map, reflect.String:
		actualLen = v.Len()
	default:
		t.Fatalf("AssertLen called on non-measurable type: %T", value)
		return
	}

	if actualLen != expectedLen {
		if len(msgAndArgs) > 0 {
			t.Errorf("Expected length %d, got %d. %v", expectedLen, actualLen, msgAndArgs[0])
		} else {
			t.Errorf("Expected length %d, got %d", expectedLen, actualLen)
		}
	}
}
