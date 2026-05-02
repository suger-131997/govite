package main

import (
	app "app"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/suger-131997/govite"
)

func main() {
	var (
		isGenMode = flag.Bool("gen", false, "generate entry points and type definitions, then exit")
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

	ctx, err := govite.WithRenderCreatorForDev(ctx, htmlTemplate, "Go + Vite Demo", viteServer, workdir)
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	fmt.Printf("Server started at http://localhost:%s\n", port)

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
