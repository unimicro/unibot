#! /bin/bash

if [ -d dist/ ] ; then
    rm -rf dist/
fi

mkdir dist/

env GOOS=linux go build -o dist/unibot.new

cp -r slack_api.token unibot.service dist/
