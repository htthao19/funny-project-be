package constant

type contextKey int

// Contexts
const (
	// ContextUID key.
	ContextUID contextKey = iota
	// ContextEmail key.
	ContextEmail
	// ContextCtx key.
	ContextCtx
)
