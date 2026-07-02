#!/usr/bin/env sh
set -eu

for file in "$@"; do
  if [ ! -f "$file" ]; then
    continue
  fi
  first_line="$(sed -n '1p' "$file")"
  if [ "$first_line" = "// SPDX-License-Identifier: MPL-2.0" ]; then
    continue
  fi
  tmp="${file}.license-tmp"
  {
    printf '%s\n' '// SPDX-License-Identifier: MPL-2.0'
    printf '%s\n' '// Copyright (c) 2026 Kernloom Contributors'
    printf '\n'
    sed -n '1,$p' "$file"
  } > "$tmp"
  mv "$tmp" "$file"
done

