#!/usr/bin/env bash

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

echo "$(date): Building container"
docker build -t kevinkamps/registrar:$VERSION .
docker tag kevinkamps/registrar:$VERSION kevinkamps/registrar:latest


echo "$(date): Pushing to docker"
echo -e "Want to push $VERSION and latest to Docker (y/n)"
read CHOICE
if [ $CHOICE = y ]
	then
	docker push kevinkamps/registrar:$VERSION
	docker push kevinkamps/registrar:latest
fi


echo "$(date): Done"