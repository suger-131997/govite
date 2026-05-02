package govite

import "context"

type titleKey struct{}

// WithTitle は title を ctx に格納し、レンダラーが HTML ページタイトルとして使用できるようにします。
// レンダーコールバックを呼び出す前に HandleFunc 内で呼び出してください。
// 設定しない場合、レンダラーはセットアップ時に指定されたデフォルトタイトルを使用します。
func WithTitle(ctx context.Context, title string) context.Context {
	return context.WithValue(ctx, titleKey{}, title)
}

// TitleFromContext は [WithTitle] によって ctx に格納されたページタイトルを取り出します。
// タイトルが設定されていない場合、第 2 戻り値は false になります。
func TitleFromContext(ctx context.Context) (string, bool) {
	value := ctx.Value(titleKey{})
	if value == nil {
		return "", false
	}
	title, ok := value.(string)
	if !ok {
		return "", false
	}
	return title, ok
}
