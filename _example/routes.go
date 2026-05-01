package app

import (
	"context"
	"net/http"
	"path"

	"app/page"
	"app/page/about"
	"app/page/about/nested"
)

func NewRouter(ctx context.Context, fsHandler http.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func() func(writer http.ResponseWriter, request *http.Request) {
		indexHandler := page.NewIndexHandler().Handler(ctx)
		notFoundHandler := page.NewNotFoundHandler().Handler(ctx)
		return func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" || r.URL.Path == "/index.html" {
				indexHandler(w, r)
				return
			}

			if path.Ext(r.URL.Path) != "" {
				fsHandler.ServeHTTP(w, r)
				return
			}

			notFoundHandler(w, r)
		}
	}())

	mux.HandleFunc("/users", page.NewUsersHandler().Handler(ctx))
	mux.HandleFunc("/users/{id}", page.NewUserHandler().Handler(ctx))
	mux.HandleFunc("/about", page.NewAboutHandler().Handler(ctx))
	mux.HandleFunc("/about/nested", about.NewIndexHandler().Handler(ctx))
	mux.HandleFunc("/about/nested/nested", nested.NewIndexHandler().Handler(ctx))

	return mux
}
