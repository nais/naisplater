naisplater
==========
![Build Status](https://github.com/nais/naisplater/workflows/Create%20docker%20image/badge.svg?branch=master)

Tool for processing go templates

# how it works

It processes a directory of go template files using a directory with yaml files as input, with the possibility to override both template files and variables for different environments.

```
$ ./naisplater --help
usage: naisplater [options] [environment] [templates_dir] [variables_dir] [output_dir] ([decryption key])

environment           specifies which subdirectory in <templates_dir> to include files from,
                      and which subdirectory in <variables_dir> to merge/override yaml with
templates_dir         directory containing go template files. Environment specific files goes into <templates_dir>/<environment>
variables_dir         directory containing yaml variable files. Environment specific overrides must go into sub dir <variables_dir>/<environment>
output_dir            folder to output processed templates (if folder exists with files, naisplater will not run)
decryption_key        secret to use for decrypting secret variables

Options:
    -h|--help         show this help
    -f|--filter       only process files matching this glob
    -n|--no-label     do not add the nais.io/created-by label
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


Example usage (whose outputs also can be in-lined into yaml files):
```
# Encrypt
cat my-credentials.json | base64 -d | openssl enc -e -aes-256-cbc -a -md md5 -A -k <your encryption key>

# Decrypt
echo -n "<encrypted string>" | openssl enc -d -aes-256-cbc -a -md md5 -A -k <your encryption key>
```

Any yaml-key with the suffix `.enc`, (see example below) will be decrypted during template interpolation.

Example:
```
mysecret.enc: U2FsdGVkX1/wy7efToqNXuQjSBYCC8F0hMBdHTQFVc0=
```

# note

- After processing the template, it will check the files for unresolved variables and error out if it finds any
- Note that variable files _must_ have same name as template file
- The environment provided as a argument is available as variable `{{ .env }}`
- Uses [tsg/gotpl](https://github.com/tsg/gotpl) for processing go templates and [mikefarah/yq](https://github.com/mikefarah/yq) for merging yaml
