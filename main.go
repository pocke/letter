package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"text/template"

	"github.com/fatih/color"
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
	commands := make(Commands, 0)
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
			c := commands[ev.GlobIndex]
			cmdStr, err := ExecTemplate(c, &TemplateArg{File: ev.Original.Name})
			if err != nil {
				panic(err)
			}

			cmd := exec.Command("bash", "-c", cmdStr)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin

			fmt.Println()
			color.New(color.Bold).Printf("Execute by letter > ")
			fmt.Println(cmdStr)

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

type Commands []*template.Template

func (c *Commands) Set(str string) error {
	t := template.New("Command").Funcs(template.FuncMap{
		"s": SubstituteForTemplate,
	})
	t, err := t.Parse(str)
	if err != nil {
		return err
	}
	*c = append(*c, t)
	return nil
}

func (c *Commands) String() string {
	return ""
}

type TemplateArg struct {
	File string
}

func SubstituteForTemplate(re, repl, src string) string {
	reg := regexp.MustCompile(re)
	return reg.ReplaceAllString(src, repl)
}

func ExecTemplate(t *template.Template, arg *TemplateArg) (string, error) {
	buf := bytes.NewBuffer([]byte{})
	err := t.Execute(buf, arg)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
