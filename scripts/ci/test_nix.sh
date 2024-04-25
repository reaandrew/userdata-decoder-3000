#!/bin/bash

# Ensure script stops on first error
set -e

# Assuming matrix variables are passed as script arguments
executable=$1
os=$2

$executable --provider aws -v --output-dir ./output/$os/
files=$(find ./output/$os -type f)
echo "Files found: $files"

# Read required files from a list file
mapfile -t requiredFiles < ./scripts/ci/expected_files.txt

# Loop through required files and check if they exist in the output
for file in "${requiredFiles[@]}"; do
  if ! grep -q "$file" <<< "$files"; then
    echo "Can't find $file"
    exit 1
  fi
done
