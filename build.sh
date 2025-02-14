#!/bin/bash
BIN_FILENAME=stargate
mkdir -p output/bin
mkdir -p output/docs
cp docs/* output/docs 2>/dev/null
cp script/* output 2>/dev/null
chmod +x output/bootstrap.sh
go build -o output/bin/${BIN_FILENAME}