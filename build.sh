#!/bin/bash
echo "Building attrezzi..."

go build -trimpath -ldflags="-w -s" -gcflags "-N -l" -o bin/att
GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-w -s" -gcflags "-N -l" -o bin/att_linux_x64