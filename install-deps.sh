#!/bin/sh

apk update
apk upgrade
apk add ca-certificates git bash make
update-ca-certificates
apk add go=1.6.2-r4

export PATH=$PATH:/usr/lib/go/bin
export BASE=`dirname $0`
#export GOPATH=/opt/rkt
#sh ${BASE}/DEPENDENCIES
exit 0

