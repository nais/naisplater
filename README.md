naisplater
==========
[![Build Status](https://github.com/nais/naisplater/workflows/Create%20docker%20image/badge.svg?branch=master)

Tool for processing go templates

# how it works

It processes a directory of go template files using a directory with yaml files as input, with the possibility to override both template files and variables for different environments.

```
$ ./naisplater --help
usage: naisplater [environment] [templates_dir] [variables_dir] [output_dir] ([decryption_key])

environment           specifies which subdirectory in <templates_dir> to include files from,
                      and which subdirectory in <variables_dir> to merge/override yaml with
templates_dir         directory containing go template files. Environment specific files goes into <templates_dir>/<environment>
variables_dir         directory containing yaml variable files. Environment specific overrides must go into sub dir <variables_dir>/<environment>
output_dir            folder to output processed templates
decryption_key        secret to use for decrypting secret variables
```

Full example (see also [test folder](https://github.com/nais/naisplater/tree/master/test))
```
$ cat templates/file      # go template file
value is {{ .value }} in environment {{ .env }}
{{ .foo }} value is {{ .value }}
$ cat templates/anotherfile  # go template file
something
$ cat templates/dev/anotherfile # override anotherfile for environment 'dev' 
overridden
$ cat vars/file       # base variable file
value: some
$ cat vars/dev/file   # variable file for environment 'dev' (will override values from base)
value: overridden
$ naisplater dev templates/ vars/ out/
-- generated file ./out/anotherfile:
something
-- generated file ./out/file:
value is overridden in environment dev
-- generated file ./out/anotherfile:
overridden
```

## encrypted variables

If you have secret variables, you can encrypt them and keep them under version control like any other variable.

Encrypt your variables like this:
```
echo -n <your secret> | openssl enc -e -aes-256-cbc -a -A -k <your encryption key>
```
The encrypted string is then put as a variable with the key-suffix `.enc`

Example:
```
mysecret.enc: U2FsdGVkX1/wy7efToqNXuQjSBYCC8F0hMBdHTQFVc0=
```

This variable will be exposed as `mysecret` during template interpolation.

# note

- After processing the template, it will check the files for unresolved variables and error out if it finds any
- Note that variable files _must_ have same name as template file
- The environment provided as a argument is available as variable `{{ .env }}`
- Uses [tsg/gotpl](https://github.com/tsg/gotpl) for processing go templates and [mikefarah/yq](https://github.com/mikefarah/yq) for merging yaml
