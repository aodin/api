package api

import (
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestMeta(t *testing.T) {
	values := url.Values{
		"campus":    {"1"},                    // Positive
		"limit":     {"10"},                   // PositiveMax
		"name":      {"x"},                    // String
		"commit":    {"true"},                 // True
		"is_active": {"false"},                // Bool
		"before":    {"2006-01-02T15:04:05Z"}, // Timestamp
		"inactive":  {""},
		"bad_uuid":  {"xxxxxx-xxxxx-xxxxx"},
		"iexact":    {"aNDroid,iphone, ,"}, // MultipleStrings
		"ids":       {"1, 2,3,"},           // MultiplePositives
	}

	meta := Meta{dirty: values, valid: url.Values{}}

	// Has
	if !meta.Has("campus") {
		t.Error("campus should be in meta")
	}
	if meta.Has("DNE") {
		t.Error("DNE should not be in meta")
	}

	// Positive
	if meta.Positive("campus") != 1 {
		t.Error("unexpected meta.Positive != 1")
	}
	if !reflect.DeepEqual(meta.valid["campus"], []string{"1"}) {
		t.Error("unexpected meta.valid post-Positive")
	}

	// // PositiveMax
	// assert.Equal(t, 0, meta.PositiveMax("limit", 5))
	// assert.Equal(t, []string(nil), meta.valid["limit"])

	// // String
	// assert.Equal(t, "x", meta.String("name"))
	// assert.Equal(t, []string{"x"}, meta.valid["name"])

	// // True
	// assert.True(t, meta.True("commit"))
	// assert.Equal(t, []string{"true"}, meta.valid["commit"])

	// Timestamp
	timestamp, err := meta.Timestamp("before")
	if err != nil {
		t.Errorf("unexpected error during meta.Timestamp: %s", err)
	}
	expected := time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)
	if !timestamp.Equal(expected) {
		t.Errorf("unexpected timestamp: %s != %s", timestamp, expected)
	}
	if !reflect.DeepEqual(
		meta.valid["before"],
		[]string{"2006-01-02T15:04:05Z"},
	) {
		t.Error("unexpected meta.valid post-Timestamp")
	}

	// // Multiple positives
	// assert.Equal(t, []int{1, 2, 3}, meta.MultiplePositives("ids"))
	// assert.Equal(t, []string{"1,2,3"}, meta.valid["ids"])

	// // Multiple strings
	// assert.Equal(t, []string{"android", "iphone"}, meta.MultipleStrings("iexact"))
	// assert.Equal(t, []string{"android,iphone"}, meta.valid["iexact"])

	// // Bool
	// b, err := meta.Bool("is_active")
	// assert.Nil(t, err, "Meta boolean parsing of is_active should not error")
	// assert.False(t, b)
	// assert.Equal(t, []string{"false"}, meta.valid["is_active"])

	// _, err = meta.Bool("inactive")
	// assert.NotNil(t, err, "Meta boolean parsing of a blank bool should error")
	// assert.Equal(t, []string(nil), meta.valid["inactive"])
}
