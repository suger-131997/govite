package govite

import "context"

type titleKey struct{}

// WithTitle stores title in ctx so that the renderer can use it as the HTML
// page title. Call this inside a HandleFunc before invoking the render
// callback. If not set, the renderer falls back to the default title provided
// at setup time.
func WithTitle(ctx context.Context, title string) context.Context {
	return context.WithValue(ctx, titleKey{}, title)
}

// TitleFromContext retrieves the page title stored in ctx by [WithTitle]. The
// second return value is false if no title was set.
func TitleFromContext(ctx context.Context) (string, bool) {
	value := ctx.Value(titleKey{})
	if value == nil {
		return "", false
	}
	title, ok := value.(string)
	if !ok {
		return "", false
	}
	return title, ok
}
