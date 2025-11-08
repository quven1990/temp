@echo off
echo "build server start..."
go env -w CGO_ENABLED=0 GOARCH=arm CC=arm-linux-gnueabihf-gcc GOOS=linux
go build  -o main main.go
echo "build server done."