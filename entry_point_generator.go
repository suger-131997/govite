package govite

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path"
	"path/filepath"
	"text/template"
)

type EntryPointGenerator struct {
	workdir  string
	rootTmpl *template.Template

	entryPoints map[string]struct{}
}

func NewEntryPointGenerator(workdir, entryPointTmpl string) *EntryPointGenerator {
	return &EntryPointGenerator{
		workdir:     workdir,
		rootTmpl:    template.Must(template.New("root").Parse(entryPointTmpl)),
		entryPoints: make(map[string]struct{}),
	}
}

type entryPointGeneratorContextKey struct{}

func WithEntryPointGenerator(ctx context.Context, generator *EntryPointGenerator) context.Context {
	return context.WithValue(ctx, entryPointGeneratorContextKey{}, generator)
}

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

func (g *EntryPointGenerator) GenerateEntryPointConfig() error {
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
