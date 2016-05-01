package api

import (
	"encoding/json"
	"io"
	"net/http"
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
	// JSON Encode is still buffered as of Go 1.6
	return json.NewEncoder(out).Encode(obj)
}

// Write writes the content type and response
func (c JSON) Write(w http.ResponseWriter, resp interface{}) error {
	w.Header().Set("Content-Type", c.MediaType())
	return c.Encode(w, resp)
}

func (c JSON) MediaType() string {
	return MediaTypeJSON
}
