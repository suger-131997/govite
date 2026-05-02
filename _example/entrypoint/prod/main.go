package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"

	"app"

	"github.com/suger-131997/govite"
)

//go:embed all:dist
var dist embed.FS

func main() {
	distFS, err := fs.Sub(dist, "dist")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	f, err := distFS.Open(".vite/manifest.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var m govite.Manifest
	err = json.NewDecoder(f).Decode(&m)
	if err != nil {
		log.Fatal(err)
	}

	ctx, err = govite.WithRenderCreatorForProd(ctx, htmlTemplate, "Go + Vite Demo", m)
	if err != nil {
		log.Fatal(err)
	}

	mux := app.NewRouter(ctx, http.FileServerFS(distFS))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	fmt.Printf("Server started at http://0.0.0.0%s\n", addr)

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}

const htmlTemplate = `
<!DOCTYPE html>
<html lang="ja">
<head>
    <meta charset="UTF-8">
    <link rel="icon" type="image/svg+xml" href="/favicon.svg" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .Title }}</title>
	<script>window.APP_PROPS={{ .AppProps }};</script>
    {{- if .StyleSheets }}
	{{ .StyleSheets }}
	{{- end }}
	{{- if .Modules }}
	{{ .Modules }}
	{{- end }}
	{{- if .PreloadModules }}
	{{ .PreloadModules }}
	{{- end }}
</head>
<body>
    <div id="root"></div>
</body>
</html>
`
