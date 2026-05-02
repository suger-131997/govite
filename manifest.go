package govite

// Manifest represents the Vite build manifest (typically .vite/manifest.json).
// The keys are the source file paths and the values are the corresponding
// [Chunk] entries produced by the build.
type Manifest map[string]*Chunk

// Chunk represents a single entry in the Vite build manifest.
type Chunk struct {
	File    string   `json:"file"`
	Name    string   `json:"name"`
	Src     string   `json:"src"`
	CSS     []string `json:"css"`
	IsEntry bool     `json:"isEntry"`
	Imports []string `json:"imports"`
}

// EntryPoint returns the entry [Chunk] whose Name matches name, or nil if no
// such entry chunk exists in the manifest.
func (m Manifest) EntryPoint(name string) *Chunk {
	for _, chunk := range m {
		if chunk.Name == name && chunk.IsEntry {
			return chunk
		}
	}
	return nil
}

// StyleSheetURLs returns the CSS file URLs for the chunk identified by name,
// including those of all transitively imported chunks. Duplicate URLs are
// omitted.
func (m Manifest) StyleSheetURLs(name string) []string {
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

// ModuleURL returns the hashed JavaScript module URL for the chunk identified
// by name. It returns an empty string if the chunk does not exist in the
// manifest.
func (m Manifest) ModuleURL(name string) string {
	chunk, ok := m[name]
	if !ok {
		return ""
	}

	return chunk.File
}

// PreloadModuleURLs returns the JavaScript module URLs suitable for
// <link rel="modulepreload"> tags for the chunk identified by name, including
// those of all transitively imported chunks. Duplicate URLs are omitted.
func (m Manifest) PreloadModuleURLs(name string) []string {
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
