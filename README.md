# govite

A Go library for integrating Go's `net/http` with React + Vite.

- Serve React pages from Go HTTP handlers
- Pass typed props from Go structs to React components
- Auto-generate TypeScript types from Go structs via reflection
- Dev mode with Vite HMR; production mode with embedded assets

## How it works

1. Define a props struct and page handler in Go
2. govite generates a React entry point and TypeScript types
3. Vite bundles the frontend (dev: HMR, prod: optimized)
4. Props are JSON-serialized and injected into the HTML as `window.APP_PROPS`

## Usage

### 1. Define a page handler

```go
package page

import (
    "context"
    "govite"
    "net/http"
)

type IndexProps struct {
    Name string `json:"name"`
}

func NewIndexHandler() *govite.PageHandler[IndexProps] {
    return govite.NewPageHandler[IndexProps](govite.PageHandlerArgs[IndexProps]{
        EntryPoint: "page/index.tsx",
        HandleFunc: func(r *http.Request, render func(ctx context.Context, props IndexProps)) {
            render(r.Context(), IndexProps{Name: "world"})
        },
    })
}
```

### 2. Write the React component

```tsx
// page/index.tsx
import type { IndexProps } from "~/types.gen.ts"

export default function IndexPage({ name }: IndexProps) {
    return <h1>Hello, {name}!</h1>
}
```

### 3. Set up dev and prod entrypoints

See [`_example/`](./_example) for a complete working example.

## Development

```bash
# Install dependencies
go mod download
npm install --prefix _example

# Start dev server (Vite on :5173, Go on :8080)
cd _example && make dev
```

## Production build

```bash
cd _example && make run
```

This will:
1. Generate `entries.gen.json` and `types.gen.ts`
2. Run `tsc -b` and `vite build`
3. Compile a self-contained Go binary with embedded assets

## Requirements

- Go 1.25+
- Node.js (npm)
- [air](https://github.com/air-verse/air) — for live reload in dev
