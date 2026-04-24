#!/usr/bin/env bash

# ─── CONFIGURATION ────────────────────────────────────────────────────────────
SOURCE_DIR="./source"
DEST_DIR="./template"
SEARCH_STRING="webtemplate"
PLACEHOLDER="{{ PROJECT_NAME }}"
# ─────────────────────────────────────────────────────────────────────────────

set -euo pipefail

if [[ ! -d "$SOURCE_DIR" ]]; then
  echo "Error: source directory '$SOURCE_DIR' does not exist." >&2
  exit 1
fi

if [[ -e "$DEST_DIR" ]]; then
  echo "Warning: destination directory '$DEST_DIR' already exists. Files will be overwritten."
fi

find "$SOURCE_DIR" -type d | while read -r dir; do
  relative="${dir#$SOURCE_DIR}"
  target_dir="${DEST_DIR}${relative}"
  mkdir -p "$target_dir"
done

find "$SOURCE_DIR" -type f | while read -r file; do
  relative="${file#$SOURCE_DIR}"
  target_file="${DEST_DIR}${relative}"

  LC_ALL=C sed "s|${SEARCH_STRING}|${PLACEHOLDER}|g" "$file" > "$target_file"
  echo "Copied: $file → $target_file"
done

ARCHIVE_NAME="$(basename "$DEST_DIR").tar.gz"
tar -czf "$ARCHIVE_NAME" -C "$(dirname "$DEST_DIR")" "$(basename "$DEST_DIR")"

echo ""
echo "Done. Files written to: $DEST_DIR"
echo "Archive created: $ARCHIVE_NAME"