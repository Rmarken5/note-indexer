package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const propertyPrefix = "  - "

type (
	Note struct {
		FileName  string
		NoteName  string
		Tags      []string
		CreatedAt time.Time
	}
)

func createNoteFromFilePath(filePath string) (Note, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Note{}, fmt.Errorf("unable to open file %s to record tags: %w", filePath, err)
	}
	defer file.Close()

	tags, err := findTagsFromFile(file)
	if err != nil {
		return Note{}, err
	}

	createdAt, err := findCreatedAtFromFile(file)
	if err != nil {
		return Note{}, err
	}

	return Note{
		FileName:  filepath.Base(filePath),
		NoteName:  strings.TrimSuffix(filepath.Base(filePath), ".md"),
		Tags:      tags,
		CreatedAt: createdAt,
	}, nil

}

func findTagsFromFile(file *os.File) ([]string, error) {
	tagFound := false
	scanner := bufio.NewScanner(file)
	tags := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		slog.Debug("Reading file", file.Name(), line)
		if !tagFound && strings.Contains(line, "tags:") {
			tagFound = true
			continue
		}
		if tagFound && strings.Contains(line, propertyPrefix) {
			tags = append(tags, strings.TrimPrefix(line, propertyPrefix))
			continue
		}
		if tagFound && !strings.Contains(line, propertyPrefix) {
			break
		}
	}
	return tags, nil
}

func findCreatedAtFromFile(file *os.File) (time.Time, error) {
	dateFound := false
	createdAt := ""

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		slog.Debug("Reading file", file.Name(), line)
		if !dateFound && strings.Contains(line, "created:") {
			dateFound = true
			continue
		}
		if dateFound && strings.Contains(line, propertyPrefix) {
			createdAt = strings.TrimPrefix(line, propertyPrefix)
			return time.Parse("2006-01-02 15:04", createdAt)
		}
	}
	stat, err := file.Stat()
	if err != nil {
		return time.Time{}, fmt.Errorf("unable to stat file %s: %w", file.Name(), err)
	}
	return stat.ModTime(), nil
}
