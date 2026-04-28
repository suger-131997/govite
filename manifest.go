package govite

type Manifest map[string]*Chunk

type Chunk struct {
	File    string   `json:"file"`
	Name    string   `json:"name"`
	Src     string   `json:"src"`
	CSS     []string `json:"css"`
	IsEntry bool     `json:"isEntry"`
	Imports []string `json:"imports"`
}

func (m Manifest) EntryPoint(name string) *Chunk {
	for _, chunk := range m {
		if chunk.Name == name && chunk.IsEntry {
			return chunk
		}
	}
	return nil
}

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

func (m Manifest) ModuleURL(name string) string {
	chunk, ok := m[name]
	if !ok {
		return ""
	}

	return chunk.File
}

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
