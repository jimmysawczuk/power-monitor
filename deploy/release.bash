#!/bin/bash

set -e

rm -rf frontend
git clone git@github.com:jimmysawczuk/power-monitor-web.git frontend
(cd frontend && yarn && PATH=toolbin:$PATH yarn build)

toolbin/go-bindata -o ./cmd/power-monitor/static.go -pkg main -prefix frontend/public frontend/public/...

BUILDTAGS="-X main.version=$(scm-status | jq -r 'if .tags | length == 0 then "" else .tags[0] end') -X main.revision=$(scm-status | jq -r '.hex.short') -X main.date=$(date --iso-8601=seconds)"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o power-monitor -ldflags "-s -d -w $BUILDTAGS" -mod=vendor ./cmd/power-monitor
