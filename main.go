package main

import (
	"os"
	"strings"

	"github.com/ogier/pflag"
)

func main() {
	err := Main()
	if err != nil {
		// TODO: error handling
		panic(err)
	}
}

func Main() error {
	w, err := NewWatcher()
	if err != nil {
		return err
	}
	defer w.Close()

	globs, commands, err := ParseFlag(os.Args)

	for _, g := range globs {
		err := w.WatchGlob(g)
		if err != nil {
			return err
		}
	}

	for {
		select {
		case ev := <-w.Event:
			logger.Println(ev)
			c := commands[ev.GlobIndex]
			if err := ExecuteCommand(c, ev.Original.Name); err != nil {
				return err
			}
		case err := <-w.Error:
			return err
		}
	}
}

func ParseFlag(args []string) ([]string, Commands, error) {
	fset := pflag.NewFlagSet(args[0], pflag.ExitOnError)
	globs := make(Strings, 0)
	commands := make(Commands, 0)
	fset.VarP(&globs, "glob", "g", "glob")
	fset.VarP(&commands, "command", "c", "command")
	fset.BoolVarP(&logger.debug, "debug", "d", false, "enable debug")

	conf, err := LoadConfig()
	if err != nil {
		return nil, nil, err
	}
	a := append(conf, args[1:]...)
	if err := fset.Parse(a); err != nil {
		return nil, nil, err
	}

	return globs, commands, nil
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
