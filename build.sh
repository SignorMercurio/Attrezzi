#!/bin/bash
echo "Building attrezzi..."

GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags="-w -s" -gcflags "-N -l" -o bin/att
GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-w -s" -gcflags "-N -l" -o bin/linux_amd64/att