package api

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/aodin/errors"
)

// Meta contains limit, offset, and optional errors for the response.
type Meta struct {
	Limit  int           `json:"limit,omitempty"`
	Offset int           `json:"offset"`
	Errors *errors.Error `json:"errors,omitempty"`
	orders Orders        // TODO keep separate?
	valid  url.Values    // successfully parsed URL parameters
	dirty  url.Values    // unsanitized URL parameters
}

// UnsanitizedGet returns the value of a key from the dirty url Values
func (meta Meta) UnsanitizedGet(key string) string {
	return meta.dirty.Get(key)
}

// DirtyGet is an alias for UnsanitizedGet
func (meta Meta) DirtyGet(key string) string {
	return meta.UnsanitizedGet(key)
}

// Delete deletes a valid key
func (meta *Meta) Delete(key string) {
	delete(meta.valid, key)
}

// Has returns true if the key exists in the meta's dirty values
func (meta Meta) Has(key string) bool {
	_, has := meta.dirty[key]
	return has
}

// Order returns the meta's ordering
func (meta Meta) Order() Orders {
	return meta.orders
}

// TODO whitelist
func (meta *Meta) ParseOrder(key string, whitelist ...Order) Orders {
	meta.orders = ParseOrders(meta.dirty.Get(key), whitelist...)
	if meta.orders.Exist() {
		meta.valid.Set(key, meta.orders.String())
	}
	return meta.orders
}

// Set sets a valid key and value
func (meta *Meta) Set(key, value string) {
	meta.valid.Set(key, value)
}

func (meta *Meta) SetOrder(orders Orders) {
	meta.orders = orders
}

// Valid returns the valid url.Values
func (meta Meta) Valid() url.Values {
	return meta.valid
}

// Sanitization
// ------------

// Bool will return the given boolean from the available URL parameters,
// saving it to the valid parameters if the boolean conversion did not fail
func (meta *Meta) Bool(key string) (b bool, err error) {
	if b, err = strconv.ParseBool(meta.dirty.Get(key)); err == nil {
		meta.valid.Set(key, strconv.FormatBool(b))
	}
	return
}

// TODO Valid formats, such as CSV, XML...
func (meta *Meta) Format() string {
	format := strings.TrimSpace(strings.ToLower(meta.dirty.Get("format")))
	if format != "" {
		meta.valid.Set("format", format)
	}
	return format
}

// I hate you sometimes, Go
type intSlice []int

func (slice intSlice) String() string {
	out := make([]string, len(slice))
	for i, n := range slice {
		out[i] = strconv.Itoa(n)
	}
	return strings.Join(out, ",")
}

// Month returns the requested month, with a possible error if invalid
func (meta *Meta) Month() (time.Month, error) {
	month := meta.Positive("month")
	if month <= 0 || month > 12 {
		return 0, fmt.Errorf("A month must be between 1 and 12 inclusive")
	}
	return time.Month(month), nil
}

// MultiplePositives will return a slice of non-zero ints from a
// given URL parameter, saving it if there is at least one.
func (meta *Meta) MultiplePositives(key string) (values []int) {
	parts := strings.Split(strings.ToLower(meta.dirty.Get(key)), ",")
	for _, part := range parts {
		if value, _ := strconv.Atoi(strings.TrimSpace(part)); value != 0 {
			values = append(values, value)
		}
	}
	if len(values) > 0 {
		meta.valid.Set(key, intSlice(values).String())
	}
	return
}

// MultipleStrings will return a slice of non-empty strings from a given
// URL parameter, saving it if there is at least one.
func (meta *Meta) MultipleStrings(key string) (values []string) {
	parts := strings.Split(strings.ToLower(meta.dirty.Get(key)), ",")
	for _, part := range parts {
		if value := strings.TrimSpace(part); value != "" {
			values = append(values, value)
		}
	}
	if len(values) > 0 {
		meta.valid.Set(key, strings.Join(values, ","))
	}
	return
}

// Positive will return the given integer from the available URL parameters,
// saving it to the valid parameters if the integer is greater than zero.
func (meta *Meta) Positive(key string) (n int) {
	if n, _ = strconv.Atoi(meta.dirty.Get(key)); n > 0 {
		meta.valid.Set(key, strconv.Itoa(n))
		return
	}
	return 0
}

// PositiveMax will return the given integer from the available URL parameters,
// saving it to the valid parameters if the integer is greater than zero
// and less than the given max.
func (meta *Meta) PositiveMax(key string, max int) int {
	n, _ := strconv.Atoi(meta.dirty.Get(key))
	if n > 0 && n <= max {
		meta.valid.Set(key, strconv.Itoa(n))
		return n
	}
	return 0
}

// String will return the given string from the available URL parameters,
// saving it to the valid parameters if it is not empty
func (meta *Meta) String(key string) (value string) {
	if value = meta.dirty.Get(key); value != "" {
		meta.valid.Set(key, value)
	}
	return
}

// Timestamp will return the given timestamp from the available URL parameters,
// saving it to the valid parameters if the timestamp parsing did not fail.
func (meta *Meta) Timestamp(key string) (timestamp time.Time, err error) {
	value := meta.dirty.Get(key)
	timestamp, err = time.Parse(time.RFC3339Nano, value)
	if err == nil {
		meta.valid.Set(key, value)
	}
	return
}

// True will return the given boolean from the available URL parameters,
// saving it to the valid parameters if the boolean is true
func (meta *Meta) True(key string) (b bool) {
	if b, _ = strconv.ParseBool(meta.dirty.Get(key)); b {
		meta.valid.Set(key, strconv.FormatBool(b))
	}
	return
}

// YearMonth returns a year and month and possibly an error
// if either year or month failed to parse
func (meta *Meta) YearMonth() (year int, month time.Month, err error) {
	// If no year or month were given, don't bother to parse
	if !(meta.Has("year") || meta.Has("month")) {
		return
	}

	month, err = meta.Month()
	if err != nil {
		return
	}

	year = meta.Positive("year") // TODO default or error values?
	return
}

func NewMeta(dirty url.Values) Meta {
	return Meta{
		valid: url.Values{},
		dirty: dirty,
	}
}
