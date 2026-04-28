package govite

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"reflect"
	"strings"
)

type renderCreatorContextKey struct{}

func WithRenderCreatorForDev(ctx context.Context, htmlTemplate, defaultTitle, viteServer, workdir string) (context.Context, error) {
	tmpl, err := template.New("index").Parse(htmlTemplate)
	if err != nil {
		return nil, err
	}

	return context.WithValue(ctx, renderCreatorContextKey{}, newDevRendererCreator(tmpl, defaultTitle, viteServer, workdir)), nil
}

func WithRenderCreatorForProd(ctx context.Context, htmlTemplate, defaultTitle string, m Manifest) (context.Context, error) {
	tmpl, err := template.New("index").Parse(htmlTemplate)
	if err != nil {
		return nil, err
	}

	return context.WithValue(ctx, renderCreatorContextKey{}, newProdRendererCreator(tmpl, defaultTitle, m)), nil
}

func RenderCreatorFromContext(ctx context.Context) (func(ctx context.Context, handler pageHandler) (Renderer, error), error) {
	value := ctx.Value(renderCreatorContextKey{})
	if value == nil {
		return nil, errors.New("no render creator found in context")
	}
	renderCreator, ok := value.(func(ctx context.Context, handler pageHandler) (Renderer, error))
	if !ok {
		return nil, errors.New("invalid render creator type")
	}

	return renderCreator, nil
}

type Renderer interface {
	Render(ctx context.Context, props any) ([]byte, error)
}

type pageHandler interface {
	EntryPoint() string
}

type devRenderer struct {
	entryPoint   string
	htmlTemplate *template.Template
	defaultTitle string

	viteServer string
	workdir    string
}

type devRendererData struct {
	AppProps   template.JS
	Title      string
	ViteServer string
	Workdir    string
	EntryPoint string
}

func newDevRendererCreator(htmlTemplate *template.Template, defaultTitle, viteServer, workdir string) func(ctx context.Context, handler pageHandler) (Renderer, error) {
	return func(ctx context.Context, handler pageHandler) (Renderer, error) {
		entryPointGenerator, err := EntryPointGeneratorFromContext(ctx)
		if err != nil {
			return nil, err
		}

		err = entryPointGenerator.Generate(handler.EntryPoint())
		if err != nil {
			return nil, err
		}

		propsTypeGenerator, err := PropsTypeGeneratorFromContext(ctx)
		if err != nil {
			return nil, err
		}

		pt := reflect.TypeOf(handler)
		m, ok := pt.MethodByName("DescribeProps")
		if !ok {
			return nil, fmt.Errorf("page handler %T does not implement DescribeProps method", handler)
		}

		propsTypeGenerator.RegisterPropsType(m.Type.In(1))

		return &devRenderer{
			entryPoint:   handler.EntryPoint(),
			htmlTemplate: htmlTemplate,
			defaultTitle: defaultTitle,
			viteServer:   viteServer,
			workdir:      workdir,
		}, nil
	}
}

func (r *devRenderer) Render(ctx context.Context, props any) ([]byte, error) {
	propsJSON, err := json.Marshal(props)
	if err != nil {
		return nil, err
	}

	title, ok := TitleFromContext(ctx)
	if !ok {
		title = "Default App Title"
	}

	data := devRendererData{
		AppProps:   template.JS(propsJSON),
		Title:      title,
		ViteServer: r.viteServer,
		Workdir:    r.workdir,
		EntryPoint: r.entryPoint,
	}

	var buf bytes.Buffer
	if err := r.htmlTemplate.Execute(&buf, data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type prodRenderer struct {
	htmlTemplate   *template.Template
	defaultTitle   string
	styleSheets    template.HTML
	modules        template.HTML
	preloadModules template.HTML
}

type prodRendererData struct {
	AppProps       template.JS
	Title          string
	StyleSheets    template.HTML
	Modules        template.HTML
	PreloadModules template.HTML
}

func newProdRendererCreator(htmlTemplate *template.Template, defaultTitle string, m Manifest) func(ctx context.Context, handler pageHandler) (Renderer, error) {
	return func(ctx context.Context, handler pageHandler) (Renderer, error) {
		chunk := m.EntryPoint(handler.EntryPoint())
		if chunk == nil {
			return nil, fmt.Errorf("entry point chunk not found: %s", handler.EntryPoint())
		}

		return &prodRenderer{
			htmlTemplate:   htmlTemplate,
			defaultTitle:   defaultTitle,
			styleSheets:    buildURLTags(`<link rel="stylesheet" href="`, `">`, m.StyleSheetURLs(chunk.Src)...),
			modules:        buildURLTags(`<script type="module" src="`, `"></script>`, m.ModuleURL(chunk.Src)),
			preloadModules: buildURLTags(`<link rel="modulepreload" href="`, `">`, m.PreloadModuleURLs(chunk.Src)...),
		}, nil
	}
}

func buildURLTags(tagPrefix, tagSuffix string, urls ...string) template.HTML {
	if len(urls) == 0 || urls[0] == "" {
		return ""
	}

	sb := strings.Builder{}
	for _, url := range urls {
		sb.WriteString(tagPrefix)
		sb.WriteString(url)
		sb.WriteString(tagSuffix)
	}
	return template.HTML(sb.String())
}

func (r *prodRenderer) Render(ctx context.Context, props any) ([]byte, error) {
	propsJSON, err := json.Marshal(props)
	if err != nil {
		return nil, err
	}

	title, ok := TitleFromContext(ctx)
	if !ok {
		title = r.defaultTitle
	}

	data := prodRendererData{
		AppProps:       template.JS(propsJSON),
		Title:          title,
		StyleSheets:    r.styleSheets,
		Modules:        r.modules,
		PreloadModules: r.preloadModules,
	}

	var buf bytes.Buffer
	if err := r.htmlTemplate.Execute(&buf, data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
