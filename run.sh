#!/bin/bash

# Build the greet example app.
run::build() {
	go build example/greet.go
}

# Install all dependencies.
run::install() {
	go get -t ./...
}

# Execute all tests.
run::test() {
	go test ./...
}

# Reformat all code.
run::format() {
	go fmt ./...
}
