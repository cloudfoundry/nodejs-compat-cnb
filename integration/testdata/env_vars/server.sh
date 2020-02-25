#!/bin/bash

source /layers/org.cloudfoundry.nodejs-compat/compat/profile.d/0_memory_available.sh

function main() {
  echo "Listening on :${PORT}"
  while true; do
    echo -e "HTTP/1.1 200 OK\n\n$(env | sort)" | nc -vv -l -p "${PORT}" -q 1
  done
}

main
