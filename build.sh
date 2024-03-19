#!/bin/bash -e
cd src
go build -o dist
echo kv options: "$@"
./dist $@
