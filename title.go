package govite

import "context"

type titleKey struct{}

func WithTitle(ctx context.Context, title string) context.Context {
	return context.WithValue(ctx, titleKey{}, title)
}

func TitleFromContext(ctx context.Context) (string, bool) {
	title, ok := ctx.Value(titleKey{}).(string)
	return title, ok
}
