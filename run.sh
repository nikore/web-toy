#!/bin/bash

# Trap sigterm and sleep before passing signal to child for $SIGNAL_TIMEOUT seconds
# this is to address a shutdown issue with traefik: https://docs.traefik.io/user-guide/marathon/#shutdown

cmd=${*}
# default to 10 seconds
SIGNAL_TIMEOUT=${SIGNAL_TIME:-10}

log() {
  echo "[$(date +%s)] run.sh -- ${1}"
}

trap_term() {
  log "Recieved sigterm, sleeping for: ${SIGNAL_TIMEOUT}s"
  sleep ${SIGNAL_TIMEOUT}s
  log "Forwarding sigterm to: ${child_pid}"
  kill -15 ${child_pid}
}

trap trap_term SIGTERM

log "About to run: ${cmd}"
${cmd} &

child_pid=$!
log "Child PID: ${child_pid}"

wait ${child_pid}
