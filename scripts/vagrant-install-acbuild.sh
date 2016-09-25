#!/bin/bash

set -xe
export DEBIAN_FRONTEND=noninteractive

pushd /vagrant
./build
sudo cp -v bin/* /usr/local/bin

curl -s -q -L -o rkt.tar.gz https://github.com/coreos/rkt/releases/download/v1.1.0/rkt-v1.1.0.tar.gz -z rkt.tar.gz
tar xfv rkt.tar.gz
sudo cp -v rkt-v1.1.0/rkt /usr/local/bin
sudo cp -v rkt-v1.1.0/*.aci /usr/local/bin
getent group rkt || sudo groupadd rkt
sudo ./rkt-v1.1.0/scripts/setup-data-dir.sh
