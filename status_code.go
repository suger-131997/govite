package govite

import "context"

type statusCodeKey struct{}

func WithStatusCode(ctx context.Context, stateCode int) context.Context {
	return context.WithValue(ctx, statusCodeKey{}, stateCode)
}

func StatusCodeFromContext(ctx context.Context) (int, bool) {
	code, ok := ctx.Value(statusCodeKey{}).(int)
	return code, ok
}
