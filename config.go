package main

import (
	"io"
	"os"

	"github.com/naoina/toml"
)

type config struct {
	Rule []rule
}

type rule struct {
	Glob    string
	Command string
}

func ConfigToArgs(r io.Reader) ([]string, error) {
	c := &config{}
	err := toml.NewDecoder(r).Decode(c)
	if err != nil {
		return nil, err
	}

	res := make([]string, 0, len(c.Rule)*4)
	for _, r := range c.Rule {
		res = append(res, "-g")
		res = append(res, r.Glob)
		res = append(res, "-c")
		res = append(res, r.Command)
	}
	return res, nil
}

func LoadConfig() ([]string, error) {
	path := "letter.toml"
	if !FileExist(path) {
		return []string{}, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return ConfigToArgs(f)
}
