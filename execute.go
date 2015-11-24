package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"syscall"
	"text/template"

	"github.com/fatih/color"
)

func ExecuteCommand(c *template.Template, fname string) error {
	cmdStr, err := ExecTemplate(c, &TemplateArg{File: fname})
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

	return nil
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
