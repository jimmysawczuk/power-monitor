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

rm -rf frontend
git clone git@github.com:jimmysawczuk/power-monitor-web.git frontend
(cd frontend && yarn && PATH=toolbin:$PATH yarn build)

toolbin/go-bindata -o ./cmd/power-monitor/static.go -pkg main -prefix frontend/public frontend/public/...

BUILDTAGS="-X main.version=$(scm-status | jq -r 'if .tags | length == 0 then "" else .tags[0] end') -X main.revision=$(scm-status | jq -r '.hex.short') -X main.date=$(date --iso-8601=seconds)"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o power-monitor -ldflags "-s -d -w $BUILDTAGS" -mod=vendor ./cmd/power-monitor
