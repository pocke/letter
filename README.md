letter
============

watch filesystem => execute command


Installation
----------

```sh
go get github.com/pocke/letter
```

Usage
--------

```sh
$ letter -g '**/*_spec.rb' -c 'rspec {{.file}}' -g 'app/**/*.rb' -c 'rspec {{.file | s "^app" "test" | s `.rb$` "_spec.rb"}}'
```
