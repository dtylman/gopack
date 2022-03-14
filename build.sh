#!/bin/bash

echo "building for linux"

go build -ldflags="-s -w"
upx gopack
