package main

import (
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
}

func NewWatcher() (*Watcher, error) {
	w, err := inotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	res := &Watcher{
		watcher: w,
		Error:   w.Error,
	}
	go res.watchEvent()

	return res, nil
}

func (w *Watcher) Close() error {
	return w.watcher.Close()
}

func (w *Watcher) WatchGlob(glob string) error {
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

func (w *Watcher) watchEvent() {
	for {
		ev := <-w.watcher.Event
		switch ev.Mask {
		case inotify.IN_MODIFY, inotify.IN_ATTRIB:
			w.Event <- &Event{
				Original: ev,
			}
		case inotify.IN_IGNORED:
			// Check file exist
			w.watcher.Watch(ev.Name)
		}
	}
}
