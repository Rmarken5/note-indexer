package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var notedDir = flag.String("dir", ".", "notes directory")

func main() {
	flag.Parse()
	notes, err := os.Open(*notedDir)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to open %s - %w", *notedDir, err))
	}
	defer notes.Close()

	err = os.Chdir(*notedDir)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to chdir to %s - %w", *notedDir, err))
	}

	indexDirFile, err := readOrCreateDirectory("Index")
	if err != nil {
		log.Fatal(err)
	}
	defer indexDirFile.Close()

	notesFiles, err := notes.ReadDir(-1)
	if err != nil {
		fmt.Println(fmt.Errorf("unable to read directory %s - %w", *notedDir, err))
	}
	noteDTOs := make([]Note, 0)
	for _, f := range notesFiles {
		if !f.IsDir() {
			note, err := createNoteFromFilePath(filepath.Join(notes.Name(), f.Name()))
			if err != nil {
				log.Fatal(err)
			}
			noteDTOs = append(noteDTOs, note)
		}
	}
	tm := TagManager{
		tagIndexDirectory: indexDirFile.Name(),
		tags:              map[string]*Tag{},
		RWMutex:           sync.RWMutex{},
	}

	for _, n := range noteDTOs {
		tm.AddTagsFromNote(n)
	}
	err = tm.WriteIndexes()
	if err != nil {
		log.Fatal(err)
	}

	tm.TitleCaseTagNames()
	err = tm.WriteMainIndex(indexDirFile.Name())
	if err != nil {
		log.Fatal(err)
	}

}

func readOrCreateDirectory(dirName string) (*os.File, error) {
	dir, err := directoryFileGetter(dirName)
	if err != nil {
		log.Printf("index directory does not exist... creating")
		err := os.Mkdir(fmt.Sprintf("%s/%s", ".", dirName), 0755)
		if err != nil {
			return nil, fmt.Errorf("unable to create index directory - %w", err)
		}
		dir, err = directoryFileGetter(dirName)
		if err != nil {
			return nil, fmt.Errorf("unable to open %s directory - %w", dirName, err)
		}
	}
	return dir, nil
}

func directoryFileGetter(dirName string) (*os.File, error) {
	return os.Open(fmt.Sprintf("%s/%s", *notedDir, dirName))
}
