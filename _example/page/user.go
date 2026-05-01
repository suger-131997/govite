package page

import (
	"context"
	"govite"
	"net/http"
	"strconv"
)

type UserDetailProps struct {
	User  User `json:"user"`
	Found bool `json:"found"`
}

func NewUserHandler() *govite.PageHandler[UserDetailProps] {
	return govite.NewPageHandler[UserDetailProps](govite.PageHandlerConfig[UserDetailProps]{
		EntryPoint: "page/user.tsx",
		HandleFunc: func(r *http.Request, render func(ctx context.Context, props UserDetailProps)) {
			id, err := strconv.Atoi(r.PathValue("id"))
			if err != nil {
				render(r.Context(), UserDetailProps{Found: false})
				return
			}

			for _, u := range dummyUsers {
				if u.ID == id {
					render(r.Context(), UserDetailProps{User: u, Found: true})
					return
				}
			}

			render(r.Context(), UserDetailProps{Found: false})
		},
	})
}
