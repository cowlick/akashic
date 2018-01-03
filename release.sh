#! /bin/bash

set -e

version=$1

executable=akashic

rm -rf release
gox -verbose -ldflags "-X main.version=${version}" -output "release/${executable}_{{.OS}}_{{.Arch}}"
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