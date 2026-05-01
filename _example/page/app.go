package page

import (
	"context"
	"govite"
	"net/http"
)

type AppProps struct{}

func NewAppHandler() *govite.PageHandler[AppProps] {
	return govite.NewPageHandler[AppProps](govite.PageHandlerConfig[AppProps]{
		EntryPoint: "page/_example.tsx",
		HandleFunc: func(r *http.Request, render func(ctx context.Context, props AppProps)) {
			render(r.Context(), AppProps{})
		},
	})
}
