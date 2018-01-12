#! /bin/bash

set -e

executable=akashic

rm -rf release
gox -verbose -output "release/${executable}_{{.OS}}_{{.Arch}}" ./$executable
cd release
for bin in *; do
  if [[ "$bin" == *windows* ]]; then
    command="${executable}.exe"
  else
    command="$executable"
  fi
  mv "$bin" "$command"
  zip "${bin}.zip" "$command"
  rm "$command"
done