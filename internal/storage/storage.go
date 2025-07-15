package storage

import "errors"

// Storage package error definitions.
var (
	// ErrUrlNotFound indicates the requested URL was not found in storage.
	// Should typically result in HTTP 404 (Not Found) response.
	ErrUrlNotFound = errors.New("url not found")

	// ErrURLExists indicates a duplicate URL/alias violation.
	ErrURLExists = errors.New("url exists")
)
