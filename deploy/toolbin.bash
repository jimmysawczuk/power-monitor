#!/bin/bash

set -e

mkdir -p toolbin
wget -q https://github.com/jimmysawczuk/tmpl/releases/download/v2.0.0/tmpl-v2.0.0-linux-amd64
mv tmpl-v2.0.0-linux-amd64 toolbin/tmpl
chmod +x toolbin/tmpl

wget -q https://github.com/jimmysawczuk/scm-status/releases/download/v2.2.0/scm-status-v2.2.0-linux-amd64
mv scm-status-v2.2.0-linux-amd64 toolbin/scm-status
chmod +x toolbin/scm-status

wget -q https://github.com/jimmysawczuk/go-bindata/releases/download/v3.1.3/go-bindata-v3.1.3-linux-amd64
mv go-bindata-v3.1.3-linux-amd64 toolbin/go-bindata
chmod +x toolbin/go-bindata
