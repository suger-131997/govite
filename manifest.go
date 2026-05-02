package govite

// Manifest は Vite のビルドマニフェスト (通常は .vite/manifest.json) を表します。
// キーはソースファイルパス、値はビルドによって生成された対応する [Chunk] です。
type Manifest map[string]*Chunk

// Chunk は Vite のビルドマニフェストの 1 エントリーを表します。
type Chunk struct {
	File    string   `json:"file"`
	Name    string   `json:"name"`
	Src     string   `json:"src"`
	CSS     []string `json:"css"`
	IsEntry bool     `json:"isEntry"`
	Imports []string `json:"imports"`
}

// EntryPoint はマニフェストの中から Name が name と一致するエントリーチャンクを返します。
// 該当するエントリーチャンクが存在しない場合は nil を返します。
func (m Manifest) EntryPoint(name string) *Chunk {
	for _, chunk := range m {
		if chunk.Name == name && chunk.IsEntry {
			return chunk
		}
	}
	return nil
}

// StyleSheets は name で識別されるチャンクの CSS ファイル URL を返します。
// 推移的にインポートされるすべてのチャンクの CSS URL も含まれます。重複する URL は除外されます。
func (m Manifest) StyleSheets(name string) []string {
	seen := make(map[string]bool)
	urls := make([]string, 0)

	var addStyleSheet func(string)
	addStyleSheet = func(name string) {
		if seen[name] {
			return
		}
		seen[name] = true

		chunk, ok := m[name]
		if !ok {
			return
		}

		urls = append(urls, chunk.CSS...)

		for _, imp := range chunk.Imports {
			addStyleSheet(imp)
		}
	}

	addStyleSheet(name)

	return urls
}

// Module は name で識別されるチャンクのハッシュ付き JavaScript モジュール URL を返します。
// マニフェストにチャンクが存在しない場合は空文字列を返します。
func (m Manifest) Module(name string) string {
	chunk, ok := m[name]
	if !ok {
		return ""
	}

	return chunk.File
}

// PreloadModules は name で識別されるチャンクに対して <link rel="modulepreload"> タグに使用する
// JavaScript モジュール URL を返します。推移的にインポートされるすべてのチャンクの URL も含まれます。
// 重複する URL は除外されます。
func (m Manifest) PreloadModules(name string) []string {
	seen := make(map[string]bool)
	urls := make([]string, 0)

	var addModulePreload func(string)
	addModulePreload = func(name string) {
		if seen[name] {
			return
		}
		seen[name] = true

		chunk, ok := m[name]
		if !ok {
			return
		}

		if chunk.File != "" {
			urls = append(urls, chunk.File)
		}

		for _, imp := range chunk.Imports {
			addModulePreload(imp)
		}
	}

	addModulePreload(name)

	return urls
}
