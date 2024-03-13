package status

const (
	// OK status.
	OK = iota + 200000
)

const (
	// Created status.
	Created = iota + 201000
)

const (
	// BadRequest error.
	BadRequest = iota + 400000
)

const (
	// Unauthorized error.
	Unauthorized = iota + 401000
)

const (
	// Forbidden error.
	Forbidden = iota + 403000
)

const (
	// NotFound error.
	NotFound = iota + 404000
)

const (
	// Conflict error.
	Conflict = iota + 409000
)

const (
	// InternalServerError error.
	InternalServerError = iota + 500000
)
