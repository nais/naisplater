#!/bin/bash -e

# 

NAIS_YAML=~/src/navikt/nais-yaml
NAISPLATER=~/src/nais/naisplater
NAIS_YAML_OUTPUT=`mktemp -d`
NAISPLATER_OUTPUT=`mktemp -d`

# don't fuck up
[ "$NAIS_YAML_OUTPUT" == "" ] && exit 1
[ "$NAISPLATER_OUTPUT" == "" ] && exit 1
[ "$NAISPLATER_DECRYPTION_KEY" == "" ] && (echo 'Set $NAISPLATER_DECRYPTION_KEY before running'; exit 1)

# build binaries
cd $NAISPLATER
mkdir -p diffs
rm -f diffs/*
go build -o bin/migrate cmd/migrate/*.go
go build -o bin/naisplater cmd/naisplater/*.go

# migrate variables
mkdir -p $NAISPLATER/vars
bin/migrate \
    --directory $NAIS_YAML/vars/ \
    --output $NAIS_YAML/new-vars/ \
    --decryption-key $NAISPLATER_DECRYPTION_KEY

for CLUSTER in knada nais-ci ci-gcp prod-gcp labs-gcp dev-gcp prod-fss dev-fss prod-sbs dev-sbs
do
    echo "Running for $CLUSTER"

    rm -rf $NAIS_YAML_OUTPUT/*
    cd $NAIS_YAML
    docker run -v `pwd`/:/app -v `pwd`/output:/utdata ghcr.io/nais/naisplater:27.0.0 /usr/bin/naisplater --no-label $CLUSTER /app/templates /app/vars /utdata $NAISPLATER_DECRYPTION_KEY

    rm -rf $NAISPLATER_OUTPUT/*
    cd $NAISPLATER

    bin/naisplater \
        --variables $NAIS_YAML/new-vars \
        --templates $NAIS_YAML/new-templates \
        --output output \
        --cluster $CLUSTER \
        --add-labels=false

    set +e
    diff -r output/ ~/src/navikt/nais-yaml/output > $NAISPLATER/diffs/$CLUSTER.diff
    set -e
done

rm -rf $NAIS_YAML_OUTPUT
rm -rf $NAISPLATER_OUTPUT
