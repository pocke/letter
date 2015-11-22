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
$ letter -g '**/*_spec.rb' -c 'rspec {{.File}}' -g 'app/**/*.rb' -c 'rspec {{.File | s "^app" "test" | s `.rb$` "_spec.rb"}}'
```
