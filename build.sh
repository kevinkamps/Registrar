#!/usr/bin/env bash

echo "$(date): Starting"
echo -e "Version: "
read VERSION

export GO111MODULE=on

mkdir bin
rm -rf bin/*

echo "$(date): Building for windows"
GOOS=windows GOARCH=386 go build -ldflags "-X main.version=$VERSION" -o ./bin/registrar_386.exe ./main.go
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$VERSION" -o ./bin/registrar_amd64.exe ./main.go

echo "$(date): Building for Linux"
GOOS=linux GOARCH=386 go build -ldflags "-X main.version=$VERSION" -o ./bin/registrar_linux_386 ./main.go
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$VERSION" -o ./bin/registrar_linux_amd64 ./main.go

echo "$(date): Building for Linux arm"
GOOS=linux GOARCH=arm go build -ldflags "-X main.version=$VERSION" -o ./bin/registrar_linux_arm ./main.go

echo "$(date): Done"