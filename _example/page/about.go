package page

import (
	"context"
	"net/http"

	"github.com/suger-131997/govite"
)

type AboutProps struct {
	Message string `json:"message"`
}

func NewAboutHandler() *govite.PageHandler[AboutProps] {
	return govite.NewPageHandler[AboutProps](govite.PageHandlerConfig[AboutProps]{
		EntryPoint: "page/about.tsx",
		HandleFunc: func(r *http.Request, render func(ctx context.Context, props AboutProps)) {
			render(r.Context(), AboutProps{
				Message: "これは About ページです。",
			})
		},
	})
}
