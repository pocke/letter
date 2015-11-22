package main

import (
	"log"
	"os"
	"path/filepath"

	"golang.org/x/exp/inotify"
)

type Event struct {
	Original *inotify.Event
}

type Watcher struct {
	watcher *inotify.Watcher
	Event   chan *Event
	Error   chan error
	globs   []string
}

func NewWatcher() (*Watcher, error) {
	w, err := inotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	res := &Watcher{
		watcher: w,
		Event:   make(chan *Event),
		Error:   w.Error,
		globs:   make([]string, 0),
	}
	go res.watchEvent()

	return res, nil
}

func (w *Watcher) Close() error {
	return w.watcher.Close()
}

func (w *Watcher) WatchGlob(glob string) error {
	w.globs = append(w.globs, glob)
	return w.watchGlob(glob)
}

func (w *Watcher) watchGlob(glob string) error {
	ms, err := filepath.Glob(glob)
	if err != nil {
		return err
	}

	for _, f := range ms {
		err := w.watcher.Watch(f)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *Watcher) Reload() error {
	err := w.Close()
	if err != nil {
		return err
	}
	wa, err := inotify.NewWatcher()
	if err != nil {
		return err
	}
	w.watcher = wa
	w.Error = wa.Error // XXX
	for _, g := range w.globs {
		err := w.watchGlob(g)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *Watcher) watchEvent() {
	for {
		ev := <-w.watcher.Event
		switch ev.Mask {
		case inotify.IN_MODIFY, inotify.IN_ATTRIB:
			w.Event <- &Event{
				Original: ev,
			}
		case inotify.IN_IGNORED:
			// TODO: Check file exist
			path := ev.Name
			if FileExist(path) {
				log.Println(ev)
				err := w.Reload()
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}

// Utility

func FileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
