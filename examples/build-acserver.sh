#!/usr/bin/env bash
set -e

echo "Building acserver..."
CGO_ENABLED=0 GOOS=linux go build -o acserver -a -tags netgo -ldflags '-w' github.com/appc/acserver

acbuild --debug begin

trap "{ export EXT=$?; acbuild --debug end && exit $EXT; }" EXIT

acbuild --debug set-name example.com/acserver
acbuild --debug copy acserver /bin/acserver
acbuild --debug copy $GOPATH/src/github.com/appc/acserver/templates /templates
acbuild --debug set-exec /bin/acserver 
acbuild --debug port add http tcp 3001
acbuild --debug mount add acis /acis
acbuild --debug label add arch amd64
acbuild --debug label add os linux
acbuild --debug write --overwrite acserver-latest-linux-amd64.aci

rm acserver
