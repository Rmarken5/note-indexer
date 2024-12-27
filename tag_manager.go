package main

import (
	_ "embed"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"text/template"
)

//go:embed templates/main_index.gotmpl
var mainIdxTmpl string

type (
	TagManager struct {
		tagIndexDirectory string
		tags              map[string]*Tag
		tagList           []*Tag
		sync.RWMutex
	}
)

func (tm *TagManager) AddTagsFromNote(note Note) {
	for _, tagString := range note.Tags {
		tm.RWMutex.Lock()
		var ok bool
		var tag *Tag
		tag, ok = tm.tags[tagString]
		if ok {
			tag.Notes = append(tag.Notes, note)
		} else {
			tm.tags[tagString] = &Tag{
				Name:  tagString,
				Notes: []Note{note},
			}
		}
		tm.RWMutex.Unlock()
	}
}

func (tm *TagManager) WriteIndexes() error {
	tm.RWMutex.RLock()
	defer tm.RWMutex.RUnlock()
	for _, t := range tm.tags {
		err := t.WriteIndex(tm.tagIndexDirectory)
		if err != nil {
			return err
		}
	}
	return nil
}

func (tm *TagManager) sortTagsByName() []*Tag {
	tags := make([]*Tag, 0, len(tm.tags))
	for _, t := range tm.tags {
		tags = append(tags, t)
	}
	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Name < tags[j].Name
	})
	tm.tagList = tags
	return tags
}

func (tm *TagManager) TitleCaseTagNames() []*Tag {
	for _, t := range tm.tags {
		t.Name = strings.ToTitle(t.Name)
	}
	return tm.sortTagsByName()
}

func (tm *TagManager) WriteMainIndex(idxDir string) error {
	f, err := os.Create(filepath.Join(idxDir, "1. Index.md"))
	if err != nil {
		return err
	}
	defer f.Close()

	parse, err := template.New("index").Parse(mainIdxTmpl)
	if err != nil {
		return err
	}
	return parse.Execute(f, tm.tagList)
}
