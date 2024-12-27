package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
)

//go:embed templates/tag_index.gotmpl
var tmpl string

type (
	Tag struct {
		Name  string
		Notes []Note
	}
)

func (t Tag) WriteIndex(idxDir string) error {
	f, err := os.Create(filepath.Join(idxDir, fmt.Sprintf("%s.md", t.Name)))
	if err != nil {
		return err
	}
	defer f.Close()

	sort.Slice(t.Notes, func(i, j int) bool {
		return t.Notes[i].CreatedAt.After(t.Notes[j].CreatedAt)
	})

	parse, err := template.New("index").Parse(tmpl)
	if err != nil {
		return err
	}
	return parse.Execute(f, t)

}
