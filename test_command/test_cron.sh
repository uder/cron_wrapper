#!/usr/bin/env bash
set -u

EXIT_CODE="$1"; shift
TIMEOUT="$1"; shift

main(){
  local exit_code="$1"; shift
  echo "STDOUT before the workload" >&1
  echo "STDERR before the workload" >&2
  sleep "$TIMEOUT"
  echo "STDOUT after the workload" >&1
  echo "STDERR after the workload" >&2
  echo "Exiting. Exit code $exit_code" >&2
  exit "$exit_code"
}

main "$EXIT_CODE"