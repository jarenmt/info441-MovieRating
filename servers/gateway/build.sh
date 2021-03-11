#!/bin/bash
GOOS=linux go build
docker build -t tilleyjaren/server .
go clean

