package govite

import (
	"context"
	"log"
	"net/http"
)

type PageHandler[T any] struct {
	entryPoint string
	handleFunc func(r *http.Request, render func(ctx context.Context, props T))
}

type PageHandlerConfig[T any] struct {
	EntryPoint string
	HandleFunc func(r *http.Request, render func(ctx context.Context, props T))
}

func NewPageHandler[T any](args PageHandlerConfig[T]) *PageHandler[T] {
	return &PageHandler[T]{
		entryPoint: args.EntryPoint,
		handleFunc: args.HandleFunc,
	}
}

func (h *PageHandler[T]) Handler(ctx context.Context) http.HandlerFunc {
	rendererCreator, err := RenderCreatorFromContext(ctx)
	if err != nil {
		panic(err) // TODO: handle error properly
	}
	renderer, err := rendererCreator(ctx, h)
	if err != nil {
		panic(err) // TODO: handle error properly
	}

	return func(w http.ResponseWriter, r *http.Request) {
		h.handleFunc(r, func(ctx context.Context, props T) {
			res, err := renderer.Render(ctx, props)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "text/html")

			if stateCode, ok := StatusCodeFromContext(ctx); ok {
				w.WriteHeader(stateCode)
			}

			if _, err := w.Write(res); err != nil {
				log.Printf("failed to write response: %v", err)
				return
			}

			return
		})
	}
}

func (h *PageHandler[T]) EntryPoint() string {
	return h.entryPoint
}

func (h *PageHandler[T]) DescribeProps(t T) {}
