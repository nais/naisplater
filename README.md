naisplater
==========

Tool for processing go templates

# how it works

It processes a directory of go template files using a directory with yaml files as input

```
$ cat files/file      # go template file
{{ .foo }} value is {{ .value }}
$ cat vars/file       # base variable file
foo: bar
value: some
$ cat vars/dev/file   # variable file for environment 'dev'
value: minimal
$ naisplater dev files/ vars/ out/
-- generated file ./out/file:
bar value is minimal 
```

- After processing the template, it will check the files for unresolved variables and error out if it finds any
- Note that variable files _must_ have same name as template file
- Uses [tsg/gotpl](https://github.com/tsg/gotpl) for processing go templates and [mikefarah/yq](https://github.com/mikefarah/yq) for merging yaml.

