# govite

Go の `net/http` と React + Vite を統合するためのライブラリです。

- Go の HTTP ハンドラから React ページを配信
- Go の構造体からリフレクションで TypeScript 型を自動生成
- 型付きのプロップスを Go から React コンポーネントへ渡す
- 開発時は Vite HMR、本番時は埋め込みアセットで配信

## 仕組み

1. Go 側でプロップスの構造体とページハンドラを定義する
2. govite がエントリポイント（TSX）と TypeScript 型定義ファイルを自動生成する
3. Vite がフロントエンドをバンドルする（開発時: HMR、本番時: 最適化済みバンドル）
4. プロップスは JSON にシリアライズされ、`window.APP_PROPS` として HTML に埋め込まれる

## インストール

```bash
go get github.com/suger-131997/govite
```

## 使い方

### 1. ページハンドラを定義する

```go
package page

import (
    "context"
    "net/http"

    "github.com/suger-131997/govite"
)

type IndexProps struct {
    Name string `json:"name"`
}

func NewIndexHandler() *govite.PageHandler[IndexProps] {
    return govite.NewPageHandler[IndexProps](govite.PageHandlerConfig[IndexProps]{
        EntryPoint: "page/index.tsx",
        HandleFunc: func(r *http.Request, render func(ctx context.Context, props IndexProps)) {
            render(r.Context(), IndexProps{Name: "世界"})
        },
    })
}
```

### 2. React コンポーネントを書く

型定義ファイルのパスは `types.gen.<パッケージパス>.d.ts` の形式で生成されます。

```tsx
// page/index.tsx
import type { IndexProps } from "~/types.gen.page.d.ts"

export default function IndexPage({ name }: IndexProps) {
    return <h1>こんにちは、{name}！</h1>
}
```

### 3. 開発用エントリポイントを作成する

開発用の `main.go` では以下を行います。

- `EntryPointGenerator` でエントリポイントの TSX ファイルを自動生成する
- `PropsTypeDefGenerator` で TypeScript 型定義ファイルを生成する
- `WithRenderCreatorForDev` で開発用レンダラーをコンテキストに設定する

```go
package main

import (
    "context"
    "log"
    "net/http"
    "os"

    "github.com/suger-131997/govite"
    "yourmodule/page"
)

func main() {
    workdir := "tmp"
    _ = os.RemoveAll(workdir)
    _ = os.MkdirAll(workdir, os.ModePerm)

    ctx := context.Background()

    entryPointGenerator := govite.NewEntryPointGenerator(workdir, entryPointTmpl)
    ctx = govite.WithEntryPointGenerator(ctx, entryPointGenerator)

    propsTypeGenerator := govite.NewPropsTypeDefGenerator()
    ctx = govite.WithPropsTypeGenerator(ctx, propsTypeGenerator)

    ctx, err := govite.WithRenderCreatorForDev(ctx, htmlTemplate, "My App", "http://localhost:5173", workdir)
    if err != nil {
        log.Fatal(err)
    }

    mux := http.NewServeMux()
    mux.HandleFunc("/", page.NewIndexHandler().Handler(ctx))

    entryPointGenerator.GenerateConfig()
    propsTypeGenerator.Generate()

    log.Fatal(http.ListenAndServe(":8080", mux))
}
```

### 4. 本番用エントリポイントを作成する

本番用の `main.go` では、Vite のビルド成果物を埋め込み (`//go:embed`) で読み込みます。

```go
package main

import (
    "context"
    "embed"
    "encoding/json"
    "io/fs"
    "log"
    "net/http"

    "github.com/suger-131997/govite"
    "yourmodule/page"
)

//go:embed all:dist
var dist embed.FS

func main() {
    distFS, _ := fs.Sub(dist, "dist")

    ctx := context.Background()

    f, _ := distFS.Open(".vite/manifest.json")
    defer f.Close()

    var m govite.Manifest
    json.NewDecoder(f).Decode(&m)

    ctx, _ = govite.WithRenderCreatorForProd(ctx, htmlTemplate, "My App", m)

    mux := http.NewServeMux()
    mux.HandleFunc("/", page.NewIndexHandler().Handler(ctx))

    log.Fatal(http.ListenAndServe(":8080", mux))
}
```

### 5. Vite の設定

[`vite-plugin-go-dev-runner`](./vite-plugin-go-dev-runner) プラグインを使うことで、`npm run dev` 実行時に Go サーバーが自動的に起動します。

```ts
// vite.config.ts
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'
import entriesConfig from './entries.gen.json'
import goDevRunner from 'vite-plugin-go-dev-runner'

export default defineConfig({
  plugins: [
    react(),
    goDevRunner({
      entry: path.resolve(__dirname, "./entrypoint/dev/main.go")
    }),
  ],
  resolve: {
    alias: { '~': path.resolve(__dirname, "./") },
  },
  build: {
    outDir: path.resolve(__dirname, "./entrypoint/prod/dist"),
    manifest: true,
    rolldownOptions: {
      input: entriesConfig,
    },
  },
})
```

## API リファレンス

### ページタイトルの設定

ハンドラ内で `govite.WithTitle` を使うと、ページごとにタイトルを変更できます。

```go
render(govite.WithTitle(r.Context(), "ページタイトル"), props)
```

### HTTP ステータスコードの設定

`govite.WithStatusCode` を使うと、レスポンスのステータスコードを指定できます。

```go
render(govite.WithStatusCode(r.Context(), http.StatusNotFound), props)
```

## サンプル

完全な動作例は [`_example/`](./_example) を参照してください。

### 開発サーバーの起動

```bash
cd _example
npm install
npm run dev
```

Vite が `:5173` で起動し、`vite-plugin-go-dev-runner` が Go サーバーを `:8080` で自動起動します。

### 本番ビルドと実行

```bash
cd _example
npm run build   # entries.gen.json と型定義を生成 → tsc → vite build → go build
npm run preview # ビルド済みバイナリを実行
```

### Cloud Run へのデプロイ

[`ko`](https://ko.build) を使ってイメージをビルドし、Cloud Run にデプロイします。

```bash
cd _example

# デフォルト設定でデプロイ（gcloud のデフォルトプロジェクトを使用）
make deploy

# プロジェクト・リージョン・サービス名・イメージリポジトリを指定する場合
make deploy \
  PROJECT_ID=my-project \
  REGION=us-central1 \
  SERVICE=my-service \
  KO_DOCKER_REPO=us-docker.pkg.dev/my-project/my-repo
```

`make deploy` は以下のステップを実行します。

1. `npm install` — npm パッケージのインストール
2. `go run entrypoint/dev/main.go -gen` — エントリポイント TSX と TypeScript 型定義を生成
3. `npx tsc -b` — TypeScript のコンパイル
4. `npx vite build` — フロントエンドアセットのビルド（`entrypoint/prod/dist/` に出力）
5. `ko build ./entrypoint/prod` — Go バイナリと埋め込みアセットをコンテナイメージとしてビルド・プッシュ
6. `gcloud run deploy` — Cloud Run へのデプロイ

## 要件

- Go 1.25+
- Node.js (npm)
- Cloud Run へのデプロイ時: [`ko`](https://ko.build)、[`gcloud` CLI](https://cloud.google.com/sdk)
