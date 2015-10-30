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
acbuild --debug set-name example.com/mongodb

# Based on ubuntu
acbuild --debug dep add quay.io/sameersbn/ubuntu

# Install mongodb
acbuild --debug run --  apt-key adv --keyserver keyserver.ubuntu.com --recv 7F0CEB10
acbuild --debug run --  /bin/sh -c 'echo "deb http://downloads-distro.mongodb.org/repo/ubuntu-upstart dist 10gen" | tee -a /etc/apt/sources.list.d/10gen.list'
acbuild --debug run --  apt-get update
acbuild --debug run --  apt-get -y install apt-utils
acbuild --debug run --  apt-get -y install mongodb-10gen

# Add a port for mongo traffic
acbuild --debug port add mongo tcp 27017

# Add a port for the mongo status page
acbuild --debug port add statuspage tcp 28017

# Run mongo
acbuild --debug set-exec -- /usr/bin/mongod --config /etc/mongodb.conf

# Write the result
acbuild --debug write --overwrite mongodb-latest-linux-amd64.aci
