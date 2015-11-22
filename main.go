package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ogier/pflag"
)

func main() {
	w, err := NewWatcher()
	if err != nil {
		panic(err)
	}
	defer w.Close()

	fset := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	globs := make(Strings, 0)
	commands := make(Strings, 0)
	fset.VarP(&globs, "glob", "g", "glob")
	fset.VarP(&commands, "command", "c", "command")

	if err := fset.Parse(os.Args[1:]); err != nil {
		panic(err)
	}

	for _, g := range globs {
		err := w.WatchGlob(g)
		if err != nil {
			panic(err)
		}
	}

	for {
		select {
		case ev := <-w.Event:
			fmt.Println(ev)
		case err := <-w.Error:
			panic(err)
		}
	}
}

// for pflag
type Strings []string

func (s *Strings) Set(str string) error {
	*s = append(*s, str)
	return nil
}

func (s *Strings) String() string {
	return strings.Join(*s, ", ")
}
