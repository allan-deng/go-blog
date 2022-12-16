#!/bin/bash

rm go-blog
rm go-blog.tar
rm -r ./tmp

echo `go version`
# CGO_ENABLED=0 GOOS=linux GOARCH=amd64 
go build -o go-blog ./main.go

mkdir tmp
mkdir -p tmp/config
cp -r ./static ./tmp
cp -r ./views ./tmp
cp ./config/*.yaml ./tmp/config
cp ./go-blog ./tmp

cd tmp
tar -cvf go-blog.tar ./*
cp go-blog.tar ..