package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
)

// JSON should implement the Serializer interface
var _ Serializer = JSON{}

func TestJSON(t *testing.T) {
	serializer := JSON{}
	// Hardcode the equality of media type in case it is overwritten
	if serializer.MediaType() != "application/json" {
		t.Errorf(
			"unexpected media type: %s != application/json",
			serializer.MediaType(),
		)
	}

	var w *httptest.ResponseRecorder

	// Response
	resp := struct {
		Message string
	}{
		Message: "yes",
	}

	// Write a response
	w = httptest.NewRecorder()
	if err := serializer.Write(w, resp); err != nil {
		t.Fatalf("valid JSON Write should not error: %s", err)
	}
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf(
			"unexpected header content type: %s != %s",
			w.Header().Get("Content-Type"), MediaTypeJSON,
		)
	}
	// http.ResponseWriter will add a newline
	b, _ := json.Marshal(resp)
	b = append(b, []byte("\n")...)
	if !bytes.Equal(w.Body.Bytes(), b) {
		t.Errorf(
			"recorded response should equal JSON output: %s != %s",
			w.Body.Bytes(), b,
		)
	}
}
