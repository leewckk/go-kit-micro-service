#!/bin/bash

SRC="go-kit-micro-service"
# DEST=`go mod why | awk 'END {print}'`
# DEST=$1

# if [ ! -n "$1" ]; then
#     DEST="github.com/leewckk/go-kit-micro-service"
# fi

rm ./go.mod 
rm ./go.sum

go mod init $DEST 
go mod tidy 
# go mod tidy -go=1.16 && go mod tidy -go=1.17

sed -i "s/$SRC/$DEST/g" `grep "$SRC" -rl --include="*.go" .`
# sed -i "s/$SRC/$DEST/g" `grep "$SRC" -rl --include="*.md" .`
# sed -i "s/$SRC/$DEST/g" `grep "$SRC" -rl ./Makefile`
# sed -i "s/$SRC/$DEST/g" `grep "$SRC" -rl ./Dockerfile`

