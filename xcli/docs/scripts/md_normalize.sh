#!/bin/bash

# ==============================================================================
# SCRIPT: normalize_md.sh
# DESCRIPTION: Normalizes any Markdown file to adhere strictly to MDRULES.md
#              (Human Readability First visual spacing constraints).
# ==============================================================================

set -euo pipefail

if [ "$#" -ne 1 ]; then
  echo "Usage: $0 <path_to_markdown_file.md>"
  exit 1
fi

INPUT_FILE="$1"

if [ ! -f "$INPUT_FILE" ]; then
  echo "Error: File '$INPUT_FILE' not found."
  exit 1
fi

# Create a temporary file for processing
TMP_FILE=$(mktemp)
trap 'rm -f "$TMP_FILE"' EXIT

# ------------------------------------------------------------------------------
# STEP 1: Read line by line and stream into state machine or direct regex
# To handle vertical spacing reliably without collapsing nested formats, we use
# a temporary buffer or an awk script which is the most robust tool for block spacing.
# ------------------------------------------------------------------------------

awk '
BEGIN {
  # Define constants
  H2_SEPARATOR = "________________________________________________________________________________"
}

# Helper function to trim trailing and leading spaces
function trim(str) {
  gsub(/^[ \t\r\n]+|[ \t\r\n]+$/, "", str)
  return str
}

{
  line = $0
  trimmed = trim(line)

  # Skip existing loose &nbsp; or separators to rebuild them cleanly
  if (trimmed == "&nbsp;" || trimmed ~ /^_____________________+$/) {
    next
  }

  # Handle H2 Headers (## Section)
  if (trimmed ~ /^## /) {
    # Check if it is not the very first line of the document to avoid leading spaces
    if (NR > 1) {
      print ""
      print ""
      print "&nbsp;"
      print "&nbsp;"
      print ""
      print ""
      print H2_SEPARATOR
      print ""
    }
    print line
    print "" # Proceeding space: exactly 1 empty line
    next
  }

  # Handle H3 Headers (### Sub-section)
  if (trimmed ~ /^### /) {
    if (NR > 1) {
      print ""
      print ""
      print "&nbsp;"
      print ""
      print ""
    }
    print line
    print "" # Proceeding space: exactly 1 empty line
    next
  }

  # Handle H4 Headers (#### Deep)
  if (trimmed ~ /^#### /) {
    if (NR > 1) {
      print ""
      print "&nbsp;"
      print ""
    }
    print line
    print "" # Proceeding space: exactly 1 empty line
    next
  }

  # Handle H5 and H6 Headers
  if (trimmed ~ /^##### / || trimmed ~ /^###### /) {
    if (NR > 1) {
      print "" # Preceding space: at least 1 empty line
    }
    print line
    # Proceeding space is optional, we preserve original stream flow
    next
  }

  # Collapse sequential multiple empty lines inside normal text to maximum of 1
  if (trimmed == "") {
    empty_count++
    if (empty_count <= 1) {
      print ""
    }
    next
  } else {
    empty_count = 0
  }

  # Print normal lines (text, code blocks, bullet points)
  print line
}
' "$INPUT_FILE" > "$TMP_FILE"

# ------------------------------------------------------------------------------
# STEP 2: Safe Overwrite
# ------------------------------------------------------------------------------
cp "$TMP_FILE" "$INPUT_FILE"

echo "Success: '$INPUT_FILE' has been normalized to Human Readability First standards."
