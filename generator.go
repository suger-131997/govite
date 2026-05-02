package govite

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"golang.org/x/mod/modfile"
)

// EntryPointGenerator はテンプレートから Vite のエントリーポイントソースファイルを生成します。
// 同一のエントリーポイントは workdir 以下に一度だけ書き込まれ、重複が防がれます。
// すべてのエントリーポイントを生成したあとに [EntryPointGenerator.GenerateConfig] を呼び出すと、
// vite.config.ts が読み込む entries.gen.json が書き出されます。
type EntryPointGenerator struct {
	workdir  string
	rootTmpl *template.Template

	entryPoints map[string]struct{}
}

// NewEntryPointGenerator は新しい EntryPointGenerator を生成します。
// エントリーポイントファイルは workdir 以下書き出されます。entryPointTmpl をテンプレートとして使用します。
func NewEntryPointGenerator(workdir, entryPointTmpl string) *EntryPointGenerator {
	return &EntryPointGenerator{
		workdir:     workdir,
		rootTmpl:    template.Must(template.New("root").Parse(entryPointTmpl)),
		entryPoints: make(map[string]struct{}),
	}
}

type entryPointGeneratorContextKey struct{}

// WithEntryPointGenerator は指定した [EntryPointGenerator] を ctx に格納し、
// 新しいコンテキストを返します。取り出すには [EntryPointGeneratorFromContext] を使用してください。
func WithEntryPointGenerator(ctx context.Context, generator *EntryPointGenerator) context.Context {
	return context.WithValue(ctx, entryPointGeneratorContextKey{}, generator)
}

// EntryPointGeneratorFromContext は [WithEntryPointGenerator] によって ctx に格納された
// [EntryPointGenerator] を取り出します。
// ジェネレーターが見つからない場合、または格納された値の型が不正な場合はエラーを返します。
func EntryPointGeneratorFromContext(ctx context.Context) (*EntryPointGenerator, error) {
	value := ctx.Value(entryPointGeneratorContextKey{})
	if value == nil {
		return nil, errors.New("entry point generator not found in context")
	}
	generator, ok := value.(*EntryPointGenerator)
	if !ok {
		return nil, errors.New("invalid entry point generator type in context")
	}
	return generator, nil
}

// Generate はジェネレーターの workdir 以下に entryPoint のファイルを書き出します。
// 同じエントリーポイントがすでに生成済みの場合、またはファイルシステム操作が失敗した場合はエラーを返します。
func (g *EntryPointGenerator) Generate(entryPoint string) error {
	if _, ok := g.entryPoints[entryPoint]; ok {
		return errors.New("entry point already exists")
	}

	g.entryPoints[entryPoint] = struct{}{}

	p := filepath.Join(g.workdir, entryPoint)

	err := os.MkdirAll(path.Dir(p), os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()

	err = g.rootTmpl.Execute(f, map[string]interface{}{
		"EntryPoint": entryPoint,
	})
	if err != nil {
		return err
	}

	return nil
}

// GenerateConfig はカレントディレクトリに entries.gen.json を書き出します。
// このファイルは登録済みエントリーポイント名と絶対パスのマッピングを持ち、
// vite.config.ts でビルド入力の設定に使用します。
func (g *EntryPointGenerator) GenerateConfig() error {
	m := make(map[string]string, len(g.entryPoints))

	for entryPoint := range g.entryPoints {
		m[entryPoint] = filepath.Join(g.workdir, entryPoint)
	}

	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile("entries.gen.json", b, 0644)
	if err != nil {
		return err
	}

	return nil
}

// PropsTypeDefGenerator は Go の [reflect.Type] を収集し、
// 対応する TypeScript の型定義ファイル (*.d.ts) を生成します。
// 型はパッケージパスでグループ化され、パッケージごとに 1 つの .d.ts ファイルが生成されます。
type PropsTypeDefGenerator struct {
	propsTypes []reflect.Type
}

// NewPropsTypeDefGenerator は空の PropsTypeDefGenerator を生成して返します。
func NewPropsTypeDefGenerator() *PropsTypeDefGenerator {
	return &PropsTypeDefGenerator{
		propsTypes: make([]reflect.Type, 0),
	}
}

type propsTypeDefGeneratorContextKey struct{}

// WithPropsTypeGenerator は指定した [PropsTypeDefGenerator] を ctx に格納し、
// 新しいコンテキストを返します。取り出すには [PropsTypeDefGeneratorFromContext] を使用してください。
func WithPropsTypeGenerator(ctx context.Context, generator *PropsTypeDefGenerator) context.Context {
	return context.WithValue(ctx, propsTypeDefGeneratorContextKey{}, generator)
}

// PropsTypeDefGeneratorFromContext は [WithPropsTypeGenerator] によって ctx に格納された
// [PropsTypeDefGenerator] を取り出します。
// ジェネレーターが見つからない場合、または格納された値の型が不正な場合はエラーを返します。
func PropsTypeDefGeneratorFromContext(ctx context.Context) (*PropsTypeDefGenerator, error) {
	value := ctx.Value(propsTypeDefGeneratorContextKey{})
	if value == nil {
		return nil, errors.New("props type generator not found in context")
	}
	generator, ok := value.(*PropsTypeDefGenerator)
	if !ok {
		return nil, errors.New("invalid props type generator type in context")
	}
	return generator, nil
}

// RegisterPropsType は [PropsTypeDefGenerator.Generate] 呼び出し時に TypeScript の型定義を
// 生成する対象として rt を登録します。
func (g *PropsTypeDefGenerator) RegisterPropsType(rt reflect.Type) {
	g.propsTypes = append(g.propsTypes, rt)
}

// Generate は登録済みのすべての props 型に対して TypeScript の型定義ファイルを書き出します。
// パッケージごとに types.gen.<pkg>.d.ts という名前のファイルが生成されます
// (例: types.gen.mypkg.mysubpkg.d.ts)。既存のファイルは上書きされます。
func (g *PropsTypeDefGenerator) Generate() error {
	grouping := make(map[string][]reflect.Type)

	for _, rt := range g.propsTypes {
		grouping[rt.PkgPath()] = append(grouping[rt.PkgPath()], rt)
	}

	modulePath, err := findModulePath()
	if err != nil {
		return err
	}

	for pkgName, rts := range grouping {
		if err := generatePropsTypeDef(strings.TrimLeft(pkgName, modulePath), rts); err != nil {
			return err
		}
	}
	return nil
}

func generatePropsTypeDef(pkgName string, propsTypes []reflect.Type) error {
	var buf bytes.Buffer
	buf.WriteString("// Code generated by govite. DO NOT EDIT.\n\n")

	generatedTypes := make(map[reflect.Type]bool)
	var typeDefinitions []string

	var processType func(t reflect.Type) string
	processType = func(t reflect.Type) string {
		switch t.Kind() {
		case reflect.String:
			return "string"
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			return "number"
		case reflect.Bool:
			return "boolean"
		case reflect.Slice, reflect.Array:
			return processType(t.Elem()) + "[]"
		case reflect.Map:
			keyType := processType(t.Key())
			valueType := processType(t.Elem())
			return "{ [key: " + keyType + "]: " + valueType + " }"
		case reflect.Struct:
			if t.PkgPath() == "time" && t.Name() == "Time" {
				return "string"
			}

			typeName := t.Name()
			if typeName == "" {
				// 匿名構造体
				var s bytes.Buffer
				s.WriteString("{\n")
				for i := range t.NumField() {
					field := t.Field(i)
					if !field.IsExported() {
						continue
					}
					tag := field.Tag.Get("json")
					fieldName := field.Name
					if tag != "" && tag != "-" {
						parts := strings.Split(tag, ",")
						fieldName = parts[0]
					}
					s.WriteString(fmt.Sprintf("  %s: %s;\n", fieldName, processType(field.Type)))
				}
				s.WriteString("}")
				return s.String()
			}

			if !generatedTypes[t] {
				generatedTypes[t] = true
				var s bytes.Buffer
				s.WriteString(fmt.Sprintf("export type %s = {\n", typeName))
				for i := range t.NumField() {
					field := t.Field(i)
					if !field.IsExported() {
						continue
					}
					tag := field.Tag.Get("json")
					fieldName := field.Name
					if tag != "" && tag != "-" {
						parts := strings.Split(tag, ",")
						fieldName = parts[0]
					}
					s.WriteString(fmt.Sprintf("  %s: %s;\n", fieldName, processType(field.Type)))
				}
				s.WriteString("}\n")
				typeDefinitions = append(typeDefinitions, s.String())
			}
			return typeName
		case reflect.Interface:
			return "any"
		default:
			return "any"
		}
	}

	for _, rt := range propsTypes {
		processType(rt)
	}

	for i := len(typeDefinitions) - 1; i >= 0; i-- {
		buf.WriteString(typeDefinitions[i])
		buf.WriteString("\n")
	}

	return os.WriteFile(fmt.Sprintf("types.gen%s.d.ts", strings.ReplaceAll(pkgName, "/", ".")), buf.Bytes(), os.ModePerm)
}

func findModulePath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		goMod := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goMod); err == nil {
			b, err := os.ReadFile(goMod)
			if err != nil {
				return "", err
			}
			f, err := modfile.Parse("go.mod", b, nil)
			if err != nil {
				return "", err
			}
			if f.Module == nil {
				return "", nil
			}
			return f.Module.Mod.Path, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", os.ErrNotExist
}
