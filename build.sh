#!/bin/bash -e
cd src
go build -o dist
echo ./dist "$@"
./dist $@
