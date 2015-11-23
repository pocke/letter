package main

import (
	"bytes"
	"testing"
)

func TestConfigToArgs(t *testing.T) {
	conf := `
[[rule]]
glob    = 'spec/**/*_spec.rb'
command = 'bundle exec rspec {{.File}}'
[[rule]]
glob    = 'app/**/*.rb'
command = 'bundle exec rspec {{.File | s "^app" "spec" | s ".rb$" "_spec.rb"}}'
`
	buf := bytes.NewBuffer([]byte(conf))
	args, err := ConfigToArgs(buf)
	if err != nil {
		t.Error(err)
	}
	t.Log(args)
	if len(args) != 8 {
		t.Errorf("len(args) should be 8, but got %d", len(args))
	}
}
