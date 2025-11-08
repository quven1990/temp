@echo off
echo "build server start..."
go env -w CGO_ENABLED=0 GOARCH=amd64 CC=arm-linux-gnueabihf-gcc GOOS=windows
go build  -o main.exe main.go
echo "build server done."