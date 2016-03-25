package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aodin/errors"
)

// Create a mock response the implements the Response interface
type response string

func (r response) AddErrors(err *errors.Error) {}

var _ Response = response("")

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

	// Write a 204
	w = httptest.NewRecorder()
	serializer.Write(w, nil, nil)
	if http.StatusNoContent != w.Code {
		t.Errorf(
			"unexpected status code: %d != %d",
			w.Code, http.StatusNoContent,
		)
	}
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf(
			"unexpected header content type: %s != %s",
			w.Header().Get("Content-Type"), MediaTypeJSON,
		)
	}

	// Write a 200
	w = httptest.NewRecorder()
	serializer.Write(w, response("OUT"), nil)
	if http.StatusOK != w.Code {
		t.Errorf(
			"unexpected status code: %d != %d",
			w.Code, http.StatusOK,
		)
	}
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf(
			"unexpected header content type: %s != %s",
			w.Header().Get("Content-Type"), MediaTypeJSON,
		)
	}

	// Write an error
	w = httptest.NewRecorder()
	serializer.Write(w, nil, errors.Meta(http.StatusBadRequest, "oops"))
	if http.StatusBadRequest != w.Code {
		t.Errorf(
			"unexpected status code: %d != %d",
			w.Code, http.StatusBadRequest,
		)
	}
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf(
			"unexpected header content type: %s != %s",
			w.Header().Get("Content-Type"), MediaTypeJSON,
		)
	}
}
