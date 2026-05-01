package main

import (
	"context"
	"flag"
	"fmt"
	"govite"
	app "govite/_example"
	"log"
	"net/http"
	"os"
)

func main() {
	var (
		isGenMode = flag.Bool("gen", false, "run in development mode")
	)
	flag.Parse()

	viteServer := "http://localhost:5173"
	workdir := "tmp"

	// Cleanup workdir
	if err := os.RemoveAll(workdir); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(workdir, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	entryPointGenerator := govite.NewEntryPointGenerator(workdir, entryPointTmpl)
	ctx = govite.WithEntryPointGenerator(ctx, entryPointGenerator)

	propsTypeGenerator := govite.NewPropsTypeDefGenerator()
	ctx = govite.WithPropsTypeGenerator(ctx, propsTypeGenerator)

	ctx, err := govite.WithRenderCreatorForDev(ctx, htmlTemplate, "govite", viteServer, workdir)
	if err != nil {
		log.Fatal(err)
	}

	mux := app.NewRouter(ctx, http.FileServerFS(os.DirFS("./public")))

	mux.Handle("/assets/", http.FileServerFS(os.DirFS(".")))

	if err := entryPointGenerator.GenerateConfig(); err != nil {
		log.Fatal(err)
	}

	if err := propsTypeGenerator.Generate(); err != nil {
		log.Fatal(err)
	}

	if *isGenMode {
		return
	}

	port := ":8080"
	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}
	fmt.Printf("Server started at http://localhost%s\n", port)
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
    <script type="module">
        import RefreshRuntime from '{{ .ViteServer }}/@react-refresh'
        RefreshRuntime.injectIntoGlobalHook(window)
        window.$RefreshReg$ = () => {}
        window.$RefreshSig$ = () => (type) => type
        window.__vite_plugin_react_preamble_installed__ = true
    </script>
    <script type="module" src="{{ .ViteServer }}/@vite/client"></script>
    <script type="module" src="{{ .ViteServer }}/{{ .Workdir }}/{{ .EntryPoint }}"></script>
</head>
<body>
    <div id="root"></div>
</body>
</html>
`

const entryPointTmpl = `
import { StrictMode } from "react";
import { createRoot } from 'react-dom/client'
import App from '~/{{ .EntryPoint }}'
import '~/global.css'


createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App {...(window.APP_PROPS || {})}/>
  </StrictMode>
)
`
