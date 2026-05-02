package govite

import "context"

type statusCodeKey struct{}

// WithStatusCode は statusCode を ctx に格納し、[PageHandler] が HTTP レスポンスのステータスコードとして
// 書き出せるようにします。レンダーコールバックを呼び出す前に HandleFunc 内で呼び出してください。
func WithStatusCode(ctx context.Context, statusCode int) context.Context {
	return context.WithValue(ctx, statusCodeKey{}, statusCode)
}

// StatusCodeFromContext は [WithStatusCode] によって ctx に格納された HTTP ステータスコードを取り出します。
// ステータスコードが設定されていない場合、第 2 戻り値は false になります。
func StatusCodeFromContext(ctx context.Context) (int, bool) {
	value := ctx.Value(statusCodeKey{})
	if value == nil {
		return 0, false
	}
	code, ok := value.(int)
	if !ok {
		return 0, false
	}

	return code, ok
}
