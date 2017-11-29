#! /bin/bash
set -e

if [ -d dist/ ] ; then
    rm -rf dist/
fi

mkdir dist/

env GOOS=linux go build -o dist/unibot.new

cp -r tokens.json unibot.service dist/
