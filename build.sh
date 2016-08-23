#!/bin/bash

mkdir -p bin
go run builder/builder.go -c config.json -o bin/horst

exit $?
