package page

import (
	"context"
	"govite"
	"net/http"
	"time"
)

type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	ZipCode string `json:"zipCode"`
}

type Tag struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type IndexProps struct {
	Name      string            `json:"name"`
	Age       int               `json:"age"`
	IsActive  bool              `json:"isActive"`
	Score     float64           `json:"score"`
	Address   Address           `json:"address"`
	Tags      []Tag             `json:"tags"`
	Metadata  map[string]string `json:"metadata"`
	CreatedAt time.Time         `json:"createdAt"`
	Website   *string           `json:"website"`
}

func NewIndexHandler() *govite.PageHandler[IndexProps] {
	return govite.NewPageHandler[IndexProps](govite.PageHandlerConfig[IndexProps]{
		EntryPoint: "page/index.tsx",
		HandleFunc: func(r *http.Request, render func(ctx context.Context, props IndexProps)) {
			website := "https://github.com/example/govite"
			render(r.Context(), IndexProps{
				Name:     "Alice",
				Age:      28,
				IsActive: true,
				Score:    92.5,
				Address: Address{
					Street:  "123 Main St",
					City:    "Tokyo",
					ZipCode: "100-0001",
				},
				Tags: []Tag{
					{Name: "go", Color: "#00ADD8"},
					{Name: "react", Color: "#61DAFB"},
					{Name: "vite", Color: "#646CFF"},
				},
				Metadata: map[string]string{
					"role":       "developer",
					"team":       "platform",
					"experience": "5 years",
				},
				CreatedAt: time.Date(2025, 1, 15, 9, 30, 0, 0, time.UTC),
				Website:   &website,
			})
		},
	})
}
