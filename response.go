package api

import "github.com/aodin/errors"

// Response must have a way to add errors
type Response interface {
	AddErrors(*errors.Error)
}
