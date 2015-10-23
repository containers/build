#!/usr/bin/env bash
set -e

acbuild --debug begin quay.io/listhub/alpine

trap "{ export EXT=$?; acbuild --debug end && exit $EXT; }" EXIT

acbuild --debug set-name example.com/nginx
acbuild --debug run apk update
acbuild --debug run apk add nginx
acbuild --debug set-exec -- /usr/sbin/nginx -g "daemon off;"
acbuild --debug port add http tcp 80
acbuild --debug mount add html /usr/share/nginx/html
acbuild --debug label add arch amd64
acbuild --debug label add os linux
acbuild --debug write --overwrite nginx-latest-linux-amd64.aci
