naisplater
==========
[![Build Naisplater](https://github.com/nais/naisplater/actions/workflows/main.yml/badge.svg)](https://github.com/nais/naisplater/actions/workflows/main.yml)

Tool for processing Go templates.

# How it works

It processes a directory of Go template files using a set of YAML variable files as input, with the possibility to override both template files and variables for different environments.

```
% naisplater --help
Usage of naisplater:
      --add-labels              add 'nais.io/created-by' and 'nais.io/touched-at' labels (default true)
      --cluster string          cluster for rendering templates and variables
      --debug                   enable debug output
      --decrypt string          decrypt all ciphertext values with 'key.enc' keys in given file; output the whole file to STDOUT
      --decryption-key string   key for decrypting variables ($NAISPLATER_DECRYPTION_KEY)
      --encrypt                 in-place encrypt all plaintext values with 'key.enc' keys
      --output string           which directory to write to
      --templates string        directory with templates
      --touched-at string       use custom timestamp in 'nais.io/touched-at' label (default "20210816T143957")
      --validate                render all templates for all clusters in-memory and check for syntax/runtime errors
      --variables string        directory with variables
```

## Building

Requires Go 1.16.

```
make
sudo install bin/naisplater /usr/local/bin
```

## Encrypted variables

If you have secret variables, you can encrypt them and keep them under version control like any other variable.

Simply modify one of the variable files, enter your secret in plain text, and give the variable an `.enc` suffix, e.g. `myvariable.enc`.
Then, run the encrypter to encrypt all unencrypted variables in-place:

```
export NAISPLATER_DECRYPTION_KEY=foo
naisplater --encrypt --variables /path/to/variables/
```

To view a file in its decrypted version, run with the `--decrypt` command pointing to a single variable file:

```
export NAISPLATER_DECRYPTION_KEY=foo
naisplater --decrypt /path/to/variables/cluster.yaml
```

Make sure unencrypted secrets are not checked in by running `git diff` before committing.

## Syntax and data validation

Run `naisplater --validate`, which will exit with non-zero status if any of the templates for any cluster fails to render for whatever reason.

```
export NAISPLATER_DECRYPTION_KEY=foo
naisplater --validate --templates /path/to/templates --variables /path/to/variables
```

# Notes

- After processing the template, it will check the files for unresolved variables and error out if it finds any
- Note that variable files _must_ have same name as the cluster they should be rendered for, or `vars.yaml` for global fallbacks
