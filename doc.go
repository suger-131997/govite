// Package govite は Go の HTTP サーバーと Vite を統合するパッケージです。
// Go からフロントエンドコンポーネントへの型安全な props の受け渡しと、
// TypeScript の型定義ファイルの自動生成をサポートします。
//
// # 概要
//
// govite は Go の HTTP バックエンドと Vite ベースのフロントエンドを橋渡しします。
// 以下の 2 つのレンダリングモードに対応しています:
//   - 開発モード: アセットのリクエストを Vite 開発サーバーにプロキシし、
//     エントリーポイントファイルをオンザフライで生成します。
//   - 本番モード: Vite のビルドマニフェストを利用してハッシュ付きアセット URL を
//     HTML レスポンスに埋め込みます。
//
// # 使い方
//
// [WithRenderCreatorForDev] または [WithRenderCreatorForProd] でレンダラーファクトリーを
// [context.Context] に付与し、[NewPageHandler] でページハンドラーを作成したうえで、
// [PageHandler.Handler] から [net/http.HandlerFunc] を取得してください。
//
// HTTP ステータスコードやページタイトルなどのページ固有の値は、
// [WithStatusCode] や [WithTitle] を使ってコンテキスト経由で渡せます。
package govite
