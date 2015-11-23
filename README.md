letter
============

[![Build Status](https://travis-ci.org/pocke/letter.svg?branch=master)](https://travis-ci.org/pocke/letter)
[![Stories in Ready](https://badge.waffle.io/pocke/letter.svg?label=ready&title=Ready)](http://waffle.io/pocke/letter)

watch filesystem => execute command


Installation
----------

```sh
go get github.com/pocke/letter
```

Usage
--------

```sh
$ letter -g '**/*_spec.rb' -c 'rspec {{.File}}' -g 'app/**/*.rb' -c 'rspec {{.File | s "^app" "test" | s `.rb$` "_spec.rb"}}'
```
