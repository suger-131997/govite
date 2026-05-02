package govite

import "context"

type statusCodeKey struct{}

// WithStatusCode stores statusCode in ctx so that [PageHandler] can write it as
// the HTTP response status code. Call this inside a HandleFunc before invoking
// the render callback.
func WithStatusCode(ctx context.Context, statusCode int) context.Context {
	return context.WithValue(ctx, statusCodeKey{}, statusCode)
}

// StatusCodeFromContext retrieves the HTTP status code stored in ctx by
// [WithStatusCode]. The second return value is false if no code was set.
func StatusCodeFromContext(ctx context.Context) (int, bool) {
	value := ctx.Value(statusCodeKey{})
	if value == nil {
		return 0, false
	}
	code, ok := value.(int)
	if !ok {
		return 0, false
	}

	return code, ok
}
