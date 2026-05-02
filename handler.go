package govite

import (
	"context"
	"log"
	"net/http"
)

// PageHandler is a generic HTTP page handler parameterized by the props type T.
// It delegates request processing to a user-supplied function and uses the
// renderer attached to the context to produce an HTML response.
type PageHandler[T any] struct {
	entryPoint string
	handleFunc func(r *http.Request, render func(ctx context.Context, props T))
}

// PageHandlerConfig holds the configuration for [NewPageHandler].
type PageHandlerConfig[T any] struct {
	// EntryPoint is the relative path to the Vite entry point file (e.g. "pages/index.ts").
	EntryPoint string
	// HandleFunc is the application logic for the page. It receives the HTTP
	// request and a render callback; calling render with a context and props
	// triggers HTML generation and writes the response.
	HandleFunc func(r *http.Request, render func(ctx context.Context, props T))
}

// NewPageHandler creates a new [PageHandler] from the given configuration.
func NewPageHandler[T any](args PageHandlerConfig[T]) *PageHandler[T] {
	return &PageHandler[T]{
		entryPoint: args.EntryPoint,
		handleFunc: args.HandleFunc,
	}
}

// Handler returns an [net/http.HandlerFunc] that processes incoming requests by
// invoking the handler's HandleFunc. It panics if a [Renderer] cannot be
// created from ctx. The renderer is created once and reused for every request.
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

// EntryPoint returns the Vite entry point path associated with this handler.
func (h *PageHandler[T]) EntryPoint() string {
	return h.entryPoint
}

// DescribeProps is a no-op method whose sole purpose is to capture the props
// type T via reflection during development-mode renderer setup. It must not be
// removed.
func (h *PageHandler[T]) DescribeProps(t T) {}
