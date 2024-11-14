#!/bin/bash

if [ -z "$1" ]; then
    echo "Error: No argument supplied"
    exit 1
fi

mkdir $1

cp -r * $1
cp -r .gitignore $1
cp -r .github $1
rm $1/create_new.sh
rm -r $1/.git || echo "No .git folder found"
rm $1/go.sum
projectName=$(basename "$1")
pushd $1 || exit 
git init
find cmd -type f | xargs -I {} sed -i "s/gotemplate/$projectName/g" {} 
find internal -type f | xargs -I {} sed -i "s/gotemplate/$projectName/g" {}
sed -i  "s/gotemplate/$projectName/g" go.mod
sed -i  "s/gotemplate/$projectName/g" Makefile
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
go install github.com/mattn/go-sqlite3

go mod tidy
make migrate
make tests

popd || exit
