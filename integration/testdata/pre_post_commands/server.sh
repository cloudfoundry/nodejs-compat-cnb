#!/bin/bash

function main() {
  echo "Listening on :${PORT}"
  while true; do
    echo -e "HTTP/1.1 200 OK\n\n$(cat package.json)" | nc -l -p "${PORT}" -q 1
  done
}

main
