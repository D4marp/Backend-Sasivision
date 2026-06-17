#!/bin/sh
set -eu

log() {
  printf '[entrypoint] %s\n' "$*"
}

run_seed() {
  if [ "${RUN_MIGRATIONS:-true}" != "true" ]; then
    log "RUN_MIGRATIONS=false — skip seed"
    return 0
  fi
  log "running migrations & seed from ${MIGRATIONS_DIR:-/app/migrations}..."
  ./seed
}

if [ "${WAIT_FOR_DB:-true}" = "true" ]; then
  attempt=0
  max="${DB_WAIT_RETRIES:-30}"
  until run_seed; do
    attempt=$((attempt + 1))
    if [ "$attempt" -ge "$max" ]; then
      log "database not ready after ${max} attempts"
      exit 1
    fi
    log "database not ready (${attempt}/${max}), retry in 2s..."
    sleep 2
  done
else
  run_seed
fi

log "starting API server on :${APP_PORT:-8080}"
export RUN_MIGRATIONS=false
exec ./server
