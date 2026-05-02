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

// EntryPointGenerator generates Vite entry point source files from a template.
// Each unique entry point is written once under workdir and tracked to avoid
// duplicates. Call [EntryPointGenerator.GenerateConfig] after all entry points
// have been generated to write the entries.gen.json file consumed by the Vite
// plugin.
type EntryPointGenerator struct {
	workdir  string
	rootTmpl *template.Template

	entryPoints map[string]struct{}
}

// NewEntryPointGenerator creates a new EntryPointGenerator that writes entry
// point files under workdir using entryPointTmpl as the file template.
// The template receives a map with a single key "EntryPoint" containing the
// relative path of the entry point.
func NewEntryPointGenerator(workdir, entryPointTmpl string) *EntryPointGenerator {
	return &EntryPointGenerator{
		workdir:     workdir,
		rootTmpl:    template.Must(template.New("root").Parse(entryPointTmpl)),
		entryPoints: make(map[string]struct{}),
	}
}

type entryPointGeneratorContextKey struct{}

// WithEntryPointGenerator stores the given [EntryPointGenerator] in ctx and
// returns the new context. Retrieve it later with
// [EntryPointGeneratorFromContext].
func WithEntryPointGenerator(ctx context.Context, generator *EntryPointGenerator) context.Context {
	return context.WithValue(ctx, entryPointGeneratorContextKey{}, generator)
}

// EntryPointGeneratorFromContext retrieves the [EntryPointGenerator] stored in
// ctx by [WithEntryPointGenerator]. It returns an error if no generator is
// found or if the stored value has an unexpected type.
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

// Generate writes an entry point file for entryPoint under the generator's
// workdir. It returns an error if the entry point was already generated or if
// a filesystem operation fails.
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

// GenerateConfig writes entries.gen.json to the current working directory.
// The file maps each registered entry point name to its absolute path and is
// read by the Vite plugin to configure build inputs.
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

// PropsTypeDefGenerator collects Go [reflect.Type] values and generates
// corresponding TypeScript type definition files (*.d.ts). Types are grouped
// by package path so that one .d.ts file is produced per package.
type PropsTypeDefGenerator struct {
	propsTypes []reflect.Type
}

// NewPropsTypeDefGenerator creates a new, empty PropsTypeDefGenerator.
func NewPropsTypeDefGenerator() *PropsTypeDefGenerator {
	return &PropsTypeDefGenerator{
		propsTypes: make([]reflect.Type, 0),
	}
}

type propsTypeDefGeneratorContextKey struct{}

// WithPropsTypeGenerator stores the given [PropsTypeDefGenerator] in ctx and
// returns the new context. Retrieve it later with
// [PropsTypeDefGeneratorFromContext].
func WithPropsTypeGenerator(ctx context.Context, generator *PropsTypeDefGenerator) context.Context {
	return context.WithValue(ctx, propsTypeDefGeneratorContextKey{}, generator)
}

// PropsTypeDefGeneratorFromContext retrieves the [PropsTypeDefGenerator] stored
// in ctx by [WithPropsTypeGenerator]. It returns an error if no generator is
// found or if the stored value has an unexpected type.
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

// RegisterPropsType adds rt to the list of types for which TypeScript
// definitions will be generated when [PropsTypeDefGenerator.Generate] is called.
func (g *PropsTypeDefGenerator) RegisterPropsType(rt reflect.Type) {
	g.propsTypes = append(g.propsTypes, rt)
}

// Generate writes TypeScript type definition files for all registered props
// types. One file per package is created with the name pattern
// types.gen.<pkg>.d.ts (e.g. types.gen.mypkg.mysubpkg.d.ts). Existing files
// are overwritten.
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
				// Anonymous struct
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
