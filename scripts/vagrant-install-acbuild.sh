#!/bin/bash

set -xe
export DEBIAN_FRONTEND=noninteractive

pushd /vagrant
./build
sudo cp -v bin/* /usr/local/bin

curl -s -q -L -o rkt.tar.gz https://github.com/coreos/rkt/releases/download/v0.9.0/rkt-v0.9.0.tar.gz -z rkt.tar.gz
tar xfv rkt.tar.gz
sudo cp -v rkt-v0.9.0/* /usr/local/bin
