# govite

Go の `net/http` と React + Vite を統合するライブラリです。

- Go の HTTP ハンドラから React ページを提供
- Go の構造体から型付き props を React コンポーネントへ渡す
- リフレクションを用いて Go の構造体から TypeScript 型定義を自動生成
- 開発時: Vite HMR、本番時: 埋め込みアセット

## 仕組み

1. Go 側で props の構造体とページハンドラを定義する
2. govite が React エントリーポイントと TypeScript 型定義を自動生成する
3. Vite がフロントエンドをバンドルする（開発時: HMR、本番時: 最適化ビルド）
4. props は JSON にシリアライズされ、`window.APP_PROPS` として HTML に埋め込まれる

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
            render(r.Context(), IndexProps{Name: "world"})
        },
    })
}
```

### 2. React コンポーネントを書く

型ファイルのパスは `types.gen.{パッケージパス}.d.ts` の形式で生成されます（例: `page` パッケージなら `types.gen.page.d.ts`）。

```tsx
// page/index.tsx
import type { IndexProps } from "~/types.gen.page.d.ts"

export default function IndexPage({ name }: IndexProps) {
    return <h1>Hello, {name}!</h1>
}
```

### 3. コンテキストのオプション

| 関数 | 説明 |
|---|---|
| `govite.WithTitle(ctx, "タイトル")` | ページの `<title>` を動的に設定する |
| `govite.WithStatusCode(ctx, 404)` | レスポンスの HTTP ステータスコードを設定する |

### 4. 開発・本番エントリーポイントのセットアップ

完全な動作例は [`_example/`](./_example) を参照してください。

## セットアップ

### 依存関係のインストール

```bash
cd _example
npm install
```

### 開発サーバーの起動

```bash
npm run dev
```

Vite サーバー（`:5173`）が起動し、付属の `vite-plugin-go-dev-runner` が Go サーバー（`:8080`）を自動的に起動・管理します。`.go` ファイルの変更を検知すると Go サーバーが自動で再起動されます。

### 本番ビルド

```bash
npm run build
```

このコマンドは以下を順番に実行します。

1. Go を `-gen` フラグ付きで実行し、`entries.gen.json` と `types.gen.*.d.ts` を生成
2. `tsc -b` で TypeScript をコンパイル
3. `vite build` でフロントエンドをバンドル
4. Go バイナリ（アセット埋め込み済み）をビルド

### 本番サーバーの起動

```bash
.bin/main
```

## 要件

- Go 1.25+
- Node.js (npm)
