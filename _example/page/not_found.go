package page

import (
	"context"
	"net/http"

	"github.com/suger-131997/govite"
)

type NotFoundProps struct {
	Path string `json:"path"`
}

func NewNotFoundHandler() *govite.PageHandler[NotFoundProps] {
	return govite.NewPageHandler[NotFoundProps](govite.PageHandlerConfig[NotFoundProps]{
		EntryPoint: "page/not_found.tsx",
		HandleFunc: func(r *http.Request, render func(ctx context.Context, props NotFoundProps)) {
			ctx := govite.WithStatusCode(r.Context(), http.StatusNotFound)
			render(ctx, NotFoundProps{Path: r.URL.Path})
		},
	})
}
