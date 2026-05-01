package app

import (
	"context"
	"govite/_example/page"
	"net/http"
)

func NewRouter(ctx context.Context, fsHandler http.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func() func(writer http.ResponseWriter, request *http.Request) {
		indexHandler := page.NewIndexHandler().Handler(ctx)
		return func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" || r.URL.Path == "/index.html" {
				indexHandler(w, r)
				return
			}

			fsHandler.ServeHTTP(w, r)
		}
	}())

	mux.HandleFunc("/users", page.NewUsersHandler().Handler(ctx))
	mux.HandleFunc("/users/{id}", page.NewUserHandler().Handler(ctx))

	return mux
}
