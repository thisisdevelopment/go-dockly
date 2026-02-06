#!/bin/bash

# Test runner script for go-dockly
# Usage: ./test-runner.sh [-v] [-c] [-r]
#   -v: verbose output
#   -c: with coverage
#   -r: with race detection

VERBOSE=""
COVERAGE=""
RACE=""

while getopts "vcr" opt; do
    case $opt in
        v) VERBOSE="-v" ;;
        c) COVERAGE="-cover" ;;
        r) RACE="-race" ;;
        *) echo "Usage: $0 [-v] [-c] [-r]"; exit 1 ;;
    esac
done

echo "Running tests..."
go test $VERBOSE $COVERAGE $RACE ./...
