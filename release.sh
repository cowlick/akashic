#! /bin/bash

set -e

commands=(
  akashic
  akashic-new
)

rm -rf release
for executable in "${commands[@]}"; do
  gox -verbose -output "release/${executable}_{{.OS}}_{{.Arch}}" ./$executable
done
cd release
for bin in akashic_*; do
  args=()
  for executable in "${commands[@]}"; do
    if [[ "$bin" == *windows* ]]; then
      command="${executable}.exe"
    else
      command="${executable}"
    fi
    args+=($command)
    mv "${bin/akashic/$executable}" "$command"
  done
  zip "${bin}.zip" ${args[*]}
  rm ${args[*]}
done
