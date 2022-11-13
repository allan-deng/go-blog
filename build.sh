#!/bin/bash

rm go-blog

echo `go version`

go build -o go-blog ./main.go
