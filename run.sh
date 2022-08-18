#!/bin/sh

# Build the greet example app.
run_build() {
	go build example/greet.go
}

# Install all dependencies
run_install() {
	go get -t ./...
}

# Execute all tests
run_test() {
	go test ./...
}

# Reformat all code
run_format() {
	go fmt ./...
}

# Static code (style) analysis
run_lint() {
  set -o errexit
  go vet ./...
  staticcheck ./...
}
