package about

import (
	"context"
	"net/http"

	"github.com/suger-131997/govite"
)

type IndexProps struct {
	Message string `json:"message"`
}

func NewIndexHandler() *govite.PageHandler[IndexProps] {
	return govite.NewPageHandler[IndexProps](govite.PageHandlerConfig[IndexProps]{
		EntryPoint: "page/about/index.tsx",
		HandleFunc: func(r *http.Request, render func(ctx context.Context, props IndexProps)) {
			render(r.Context(), IndexProps{
				Message: "これは About のネストされたページです。",
			})
		},
	})
}
