package govite

import "context"

type statusCodeKey struct{}

// WithStatusCode stores statusCode in ctx so that [PageHandler] can write it as
// the HTTP response status code. Call this inside a HandleFunc before invoking
// the render callback.
func WithStatusCode(ctx context.Context, stateCode int) context.Context {
	return context.WithValue(ctx, statusCodeKey{}, stateCode)
}

// StatusCodeFromContext retrieves the HTTP status code stored in ctx by
// [WithStatusCode]. The second return value is false if no code was set.
func StatusCodeFromContext(ctx context.Context) (int, bool) {
	code, ok := ctx.Value(statusCodeKey{}).(int)
	return code, ok
}
