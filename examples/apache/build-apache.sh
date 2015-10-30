#!/usr/bin/env bash
set -e

if [ "$EUID" -ne 0 ]; then
    echo "This script uses functionality which requires root privileges"
    exit 1
fi

# Start the build with an empty ACI
acbuild --debug begin

# In the event of the script exiting, end the build
trap "{ export EXT=$?; acbuild --debug end && exit $EXT; }" EXIT

# Name the ACI
acbuild --debug set-name example.com/apache

# Based on alpine
acbuild --debug dep add quay.io/coreos/alpine-sh

# Install apache
acbuild --debug run -- apk update
acbuild --debug run -- apk add apache2

acbuild --debug run -- /bin/sh -c "echo 'ServerName localhost' >> /etc/apache2/httpd.conf"

# Add a port for http traffic on port 80
acbuild --debug port add http tcp 80

# Add a mount point for files to serve
acbuild --debug mount add html /var/www/localhost/htdocs

# Run apache, and remain in the foreground
acbuild --debug set-exec -- /bin/sh -c "chmod 755 / && /usr/sbin/httpd -D FOREGROUND"

# Write the result
acbuild --debug write --overwrite apache-latest-linux-amd64.aci
