package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/aodin/errors"
)

const MediaTypeJSON = "application/json"

// JSON implements JSON serialization
type JSON struct{}

var _ Serializer = JSON{}

// Decode decodes the given ReadCloser into the given JSON destination
func (c JSON) Decode(data io.ReadCloser, dest interface{}) error {
	defer data.Close()
	return json.NewDecoder(data).Decode(dest)
}

// Encode writes the given interface as JSON
func (c JSON) Encode(out io.Writer, obj interface{}) error {
	// Encoding is still buffered as of Go 1.6
	return json.NewEncoder(out).Encode(obj)
}

// Write writes the response code, content type and response
func (c JSON) Write(w http.ResponseWriter, r interface{}, err *errors.Error) {
	w.Header().Set("Content-Type", c.MediaType())

	if err == nil && r == nil {
		// No errors and an empty response: status 204
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil { // TODO Confirm that the err is empty?
		// Error response
		w.WriteHeader(err.Code)
		c.Encode(w, r)
		return
	}

	// Otherwise, write a status OK response
	c.Encode(w, r)
}

func (c JSON) MediaType() string {
	return MediaTypeJSON
}
