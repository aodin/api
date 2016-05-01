package api

import (
	"io"
	"net/http"
)

// Serializer reads and writes API responses
type Serializer interface {
	Decoder
	Encoder
	MediaType() string
	Write(http.ResponseWriter, interface{}) error
}

// Decoder is the common decoding interface
type Decoder interface {
	Decode(io.ReadCloser, interface{}) error
}

// Encoder is the common encoding interface
type Encoder interface {
	Encode(io.Writer, interface{}) error
}
