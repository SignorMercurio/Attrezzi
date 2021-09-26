#!/bin/bash
echo "Building attrezzi..."

CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags="-w -s" -gcflags "-N -l" -o bin/att
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-w -s" -gcflags "-N -l" -o bin/linux_amd64/att