package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"govite"
	app "govite/_example"
	"io/fs"
	"log"
	"net/http"
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

	ctx, err = govite.WithRenderCreatorForProd(ctx, htmlTemplate, "govite", m)
	if err != nil {
		log.Fatal(err)
	}

	mux := app.NewRouter(ctx, http.FileServerFS(distFS))

	port := ":8080"
	fmt.Printf("Server started at http://localhost%s\n", port)

	server := &http.Server{
		Addr:    "localhost:8080",
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
