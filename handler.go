package govite

import (
	"context"
	"log"
	"net/http"
)

// PageHandler は props の型 T をパラメーターとするジェネリックな HTTP ページハンドラーです。
// リクエスト処理をユーザーが提供する関数に委譲し、コンテキストに付与されたレンダラーを使って
// HTML レスポンスを生成します。
type PageHandler[T any] struct {
	entryPoint string
	handleFunc func(r *http.Request, render func(ctx context.Context, props T))
}

// PageHandlerConfig は [NewPageHandler] の設定を保持します。
type PageHandlerConfig[T any] struct {
	// EntryPoint は Vite のエントリーポイントファイルへの相対パスです (例: "pages/index.ts")。
	EntryPoint string
	// HandleFunc はページのアプリケーションロジックです。HTTP リクエストとレンダーコールバックを受け取り、
	// コンテキストと props を渡してレンダーコールバックを呼び出すと HTML 生成とレスポンス書き込みが行われます。
	HandleFunc func(r *http.Request, render func(ctx context.Context, props T))
}

// NewPageHandler は指定した設定から新しい [PageHandler] を生成して返します。
func NewPageHandler[T any](args PageHandlerConfig[T]) *PageHandler[T] {
	return &PageHandler[T]{
		entryPoint: args.EntryPoint,
		handleFunc: args.HandleFunc,
	}
}

// Handler はハンドラーの HandleFunc を呼び出してリクエストを処理する [net/http.HandlerFunc] を返します。
// ctx から [Renderer] を生成できない場合はパニックします。
// レンダラーは一度だけ生成され、以降のリクエストで再利用されます。
func (h *PageHandler[T]) Handler(ctx context.Context) http.HandlerFunc {
	var renderer Renderer
	var setupErr error
	if rendererCreator, err := RenderCreatorFromContext(ctx); err != nil {
		setupErr = err
	} else if renderer, err = rendererCreator(ctx, h); err != nil {
		setupErr = err
	}

	return func(w http.ResponseWriter, r *http.Request) {
		h.handleFunc(r, func(ctx context.Context, props T) {
			if setupErr != nil {
				http.Error(w, setupErr.Error(), http.StatusInternalServerError)
				return
			}

			res, err := renderer.Render(ctx, props)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "text/html")

			if stateCode, ok := StatusCodeFromContext(ctx); ok {
				w.WriteHeader(stateCode)
			}

			if _, err := w.Write(res); err != nil {
				log.Printf("failed to write response: %v", err)
				return
			}

			return
		})
	}
}

// EntryPoint はこのハンドラーに関連付けられた Vite のエントリーポイントパスを返します。
func (h *PageHandler[T]) EntryPoint() string {
	return h.entryPoint
}

// DescribeProps は何も行わないメソッドです。開発モードのレンダラーセットアップ時に
// リフレクションで props の型 T を取得するためだけに存在します。削除しないでください。
func (h *PageHandler[T]) DescribeProps(t T) {}
