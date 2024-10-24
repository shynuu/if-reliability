#!/bin/bash

# Usage: ./compile.sh <source-file>
# This script compiles the specified Go source file using `go build`.
# It checks if the file exists and if there are any changes since the last build.
# If the file does not exist, it exits with an error message.
# If no changes are detected and the compiled binary already exists, it skips compilation.
# Otherwise, it builds the binary named 'if-reliability' from the given source file.

compile() {
  if [ ! -f "$1" ]; then
    echo "File $1 does not exist"
    exit 1
  fi

  if [ -f "if-reliability" ] && [ "$1" -nt "if-reliability" ]; then
    echo "No changes detected. Skipping compilation."
    exit 0
  fi

  go build -o if-reliability "$1"
}

compile $1