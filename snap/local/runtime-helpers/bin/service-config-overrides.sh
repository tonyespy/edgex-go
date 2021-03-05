#!/bin/bash -e

# convert cmdline to string array
ARGV=($@)

# grab binary path
BINPATH="${ARGV[0]}"
logger "edgex service override: BINPATH=$BINPATH"

# binary name == service name/key
SERVICE=$(basename "$BINPATH")
logger "edgex service override: SERVICE=$SERVICE"

SERVICE_ENV="$SNAP_DATA/config/$SERVICE/res/$SERVICE.env"
logger "edgex service override: : SERVICE_ENV=$SERVICE_ENV"

if [ -f "$SERVICE_ENV" ]; then
    logger "edgex service override: : sourcing $SERVICE_ENV"
    source "$SERVICE_ENV"
fi

exec "$@"
