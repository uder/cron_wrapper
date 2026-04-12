#!/usr/bin/env bash
set -u

EXIT_CODE="$1"; shift
TIMEOUT="$1"; shift

main(){
  local exit_code="$1"; shift
  echo "StdOut before the workload" >&1
  echo "StdErr before the workload" >&2
  sleep "$TIMEOUT"
  echo "StdOut after the workload" >&1
  echo "StdErr after the workload" >&2
  echo "Exiting. Exit code $exit_code" >&2
  exit "$exit_code"
}

main "$EXIT_CODE"
