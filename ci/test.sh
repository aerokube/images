#!/bin/bash

export GO111MODULE="on"
go test -race -v -coverprofile=coverage.txt -covermode=atomic ./...
