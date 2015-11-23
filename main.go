package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
	"text/template"

	"github.com/fatih/color"
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
			cmdStr, err := ExecTemplate(c, &TemplateArg{File: ev.Original.Name})
			if err != nil {
				return err
			}

			cmd := exec.Command("bash", "-c", cmdStr)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin

			fmt.Println()
			color.New(color.Bold).Printf("Execute by letter > ")
			fmt.Println(cmdStr)

			err = cmd.Run()
			status, err := ExitStatus(err)
			if err != nil {
				return err
			}
			color.New(color.Bold).Print("Finished command with status code ")

			var colorAttr color.Attribute
			if status == 0 {
				colorAttr = color.FgGreen
			} else {
				colorAttr = color.FgRed
			}
			color.New(color.Bold, colorAttr).Println(status)
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

func ExitStatus(err error) (code int, e error) {
	defer func() {
		if r := recover(); r != nil {
			e = fmt.Errorf("%+v", r)
		}
	}()

	if err == nil {
		return 0, nil
	}

	return err.(*exec.ExitError).Sys().(syscall.WaitStatus).ExitStatus(), nil
}
