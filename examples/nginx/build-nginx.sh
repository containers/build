#!/usr/bin/env bash
set -e

if [ "$EUID" -ne 0 ]; then
    echo "This script uses functionality which requires root privileges"
    exit 1
fi

# Start the build with an empty ACI
acbuild --debug begin

# In the event of the script exiting, end the build
acbuildEnd() {
    export EXIT=$?
    acbuild --debug end && exit $EXIT 
}
trap acbuildEnd EXIT

# Name the ACI
acbuild --debug set-name example.com/nginx

# Based on alpine
acbuild --debug dep add quay.io/coreos/alpine-sh

# Install nginx
acbuild --debug run apk update
acbuild --debug run apk add nginx

# Add a port for http traffic over port 80
acbuild --debug port add http tcp 80

# Add a mount point for files to serve
acbuild --debug mount add html /usr/share/nginx/html

# Run nginx in the foreground
acbuild --debug set-exec -- /usr/sbin/nginx -g "daemon off;"

# Save the ACI
acbuild --debug write --overwrite nginx-latest-linux-amd64.aci
