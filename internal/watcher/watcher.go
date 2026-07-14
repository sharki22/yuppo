package watcher

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
)

type Watcher struct {
	fsWatcher *fsnotify.Watcher
	rootPath  string
	matcher   gitignore.Matcher
	onChange  func()
}

func New(rootPath string, matcher gitignore.Matcher, onChange func()) (*Watcher, error) {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("create fsnotify watcher: %w", err)
	}

	w := &Watcher{
		fsWatcher: fsw,
		rootPath:  rootPath,
		matcher:   matcher,
		onChange:  onChange,
	}

	if err := w.scan(rootPath); err != nil {
		fsw.Close()
		return nil, fmt.Errorf("scan directories: %w", err)
	}

	go w.loop()

	return w, nil
}

func (w *Watcher) Close() error {
	return w.fsWatcher.Close()
}

func (w *Watcher) scan(dir string) error {
	if err := w.fsWatcher.Add(dir); err != nil {
		return fmt.Errorf("add watch %s: %w", dir, err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read dir %s: %w", dir, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if entry.Name() == ".git" {
			continue
		}

		fullPath := filepath.Join(dir, entry.Name())

		relPath, err := filepath.Rel(w.rootPath, fullPath)
		if err != nil {
			continue
		}
		parts := strings.Split(relPath, string(filepath.Separator))
		if w.matcher.Match(parts, true) {
			continue
		}

		if err := w.scan(fullPath); err != nil {
			log.Printf("warning: skip %s: %v", fullPath, err)
		}
	}

	return nil
}

func (w *Watcher) loop() {
	for {
		select {
		case event, ok := <-w.fsWatcher.Events:
			if !ok {
				return
			}
			w.handleEvent(event)
		case err, ok := <-w.fsWatcher.Errors:
			if !ok {
				return
			}
			log.Printf("watcher error: %v", err)
		}
	}
}

func (w *Watcher) handleEvent(event fsnotify.Event) {
	relPath, err := filepath.Rel(w.rootPath, event.Name)
	if err != nil {
		return
	}

	parts := strings.Split(relPath, string(filepath.Separator))

	isDir := false
	if event.Has(fsnotify.Create) {
		if fi, err := os.Stat(event.Name); err == nil && fi.IsDir() {
			isDir = true
		}
	}

	if w.matcher.Match(parts, isDir) {
		return
	}

	switch {
	case event.Has(fsnotify.Create) && isDir:
		if filepath.Base(event.Name) == ".git" {
			return
		}
		if err := w.fsWatcher.Add(event.Name); err != nil {
			log.Printf("warning: watch new dir %s: %v", event.Name, err)
			return
		}
		w.onChange()

	case event.Has(fsnotify.Create) && !isDir:
		w.onChange()

	case event.Has(fsnotify.Write) && !isDir:
		w.onChange()

	case event.Has(fsnotify.Remove):
		w.onChange()
	}
}
