#!/bin/bash -e
cd src
go build -o dist
./dist $@
