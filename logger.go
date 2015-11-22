package main

import "log"

type Logger struct {
	debug bool
}

var logger = &Logger{}

func (l *Logger) Println(v ...interface{}) {
	if l.debug {
		log.Println(v...)
	}
}

func init() {
	log.SetFlags(log.Lshortfile)
}
