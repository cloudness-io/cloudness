#!/usr/bin/env sh

echo "Updating cmd/app/wire_gen.go"
go tool wire gen github.com/cloudness-io/cloudness/cmd/app
# format generated file as we can't exclude it from being formatted easily.
go tool goimports -w ./cmd/app/wire_gen.go