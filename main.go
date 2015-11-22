package main

import (
	"os"
	"os/exec"
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
	fset.BoolVarP(&logger.debug, "debug", "d", false, "enable debug")

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
			logger.Println(ev)
			c := strings.Split(commands[ev.GlobIndex], " ")
			cmd := exec.Command(c[0], c[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			cmd.Run()
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
