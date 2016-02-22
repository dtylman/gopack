#!/usr/bin/env bash

echo "gofmt..."
gofmt -w ./.
echo "goimports..."
goimports -w ./.
echo "building..."
go build ./...
go build
echo "testing..."
go test ./...

echo "vet..."
go tool vet -composites=false ./.
echo "lint.."
golint ./...

echo "done"