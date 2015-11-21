package main

import "golang.org/x/exp/inotify"

type Event struct {
}

type Watcher struct {
	watcher *inotify.Watcher
	Event   chan *Event
}

func NewWatcher() (*Watcher, error) {
	w, err := inotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	res := &Watcher{
		watcher: w,
	}
	return res, nil
}

func (w *Watcher) Close() error {
	return w.watcher.Close()
}

func (w *Watcher) Error() chan error {
	return w.watcher.Error
}

func WatchGlob(glob string) error {
	// TODO: implement
	return nil
}

func (w *Watcher) watchEvent() {
	// TODO: implement
}
