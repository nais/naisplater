naisplater
==========

Tool for processing go templates

# how it works

It processes a directory of go template files using a directory with yaml files as input

```
$ ./naisplater --help
usage: naisplater [environment] [templates_dir] [variables_dir] [output_dir]

environment           specifies which subdirectory in <templates_dir> to include files from,
                      and which subdirectory in <variables_dir> to merge/override yaml with
templates_dir         directory containing go template files. Environment specific files goes into <templates_dir>/<environment>
variables_dir         directory containing yaml variable files. Environment specific overrides must go into sub dir <variables_dir>/<environment>
output_dir            folder to output processed templates
```

Full example:
```
$ cat templates/file      # go template file
{{ .foo }} value is {{ .value }}
$ cat templates/anotherfile  # go template file
something
$ cat templates/dev/anotherfile # override anotherfile for environment 'dev' 
overridden
$ cat vars/file       # base variable file
foo: bar
value: some
$ cat vars/dev/file   # variable file for environment 'dev' (will override values from base)
value: minimal
$ naisplater dev templates/ vars/ out/
-- generated file ./out/anotherfile:
something
-- generated file ./out/file:
bar value is minimal 
-- generated file ./out/anotherfile:
overridden
```

- After processing the template, it will check the files for unresolved variables and error out if it finds any
- Note that variable files _must_ have same name as template file
- Uses [tsg/gotpl](https://github.com/tsg/gotpl) for processing go templates and [mikefarah/yq](https://github.com/mikefarah/yq) for merging yaml

