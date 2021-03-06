package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/bmatcuk/doublestar"

	"golang.org/x/exp/inotify"
)

type Event struct {
	Original  *inotify.Event
	GlobIndex int
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
		Error:   make(chan error),
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
	ms, err := doublestar.Glob(glob)
	if err != nil {
		return err
	}

	for _, f := range ms {
		err := w.watcher.Watch(f)
		if err != nil {
			return err
		}
		logger.Println("Watching ", f, " ...")
	}
	return nil
}

func (w *Watcher) Reload() error {
	logger.Println("Reloading watcher...")
	err := w.Close()
	if err != nil {
		return err
	}

	wa, err := inotify.NewWatcher()
	if err != nil {
		return err
	}
	w.watcher = wa

	for _, g := range w.globs {
		err := w.watchGlob(g)
		if err != nil {
			return err
		}
	}

	go w.watchEvent()
	return nil
}

func (w *Watcher) watchEvent() {
	for {
		select {
		case ev := <-w.watcher.Event:
			if ev == nil {
				return
			}
			switch ev.Mask {
			case inotify.IN_MODIFY, inotify.IN_ATTRIB:
				w.onModify(ev)
			case inotify.IN_IGNORED:
				if w.onDelete(ev) {
					return //When reloaded
				}
			}
		case err := <-w.watcher.Error:
			if err != nil {
				w.Error <- err
			}
			return
		}
	}
}

func (w *Watcher) onModify(ev *inotify.Event) {
	logger.Println("Event: ", ev)

	var idx int
	for i, g := range w.globs {
		if ok, _ := filepath.Match(g, ev.Name); ok {
			idx = i
			break
		}
	}
	w.Event <- &Event{
		Original:  ev,
		GlobIndex: idx,
	}
}

// return true when reloaded
func (w *Watcher) onDelete(ev *inotify.Event) bool {
	logger.Println("Event: ", ev)

	path := ev.Name
	time.Sleep(10 * time.Millisecond)
	if !FileExist(path) {
		return false
	}

	err := w.Reload()
	if err != nil {
		logger.Println(err)
	}
	return true
}

// Utility

func FileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
