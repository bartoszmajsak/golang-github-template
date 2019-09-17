#!/bin/bash

die () {
    echo >&2 "$@"
    exit 1
}

defaultName=$(git remote -v | cut -d':' -f 2 | cut -d'.' -f 1 | uniq)
read -p "Repo name (${defaultName}): " name
test -z "$name" && name=${defaultName}

find . -type f -not -path "./vendor/*" -not -path "./.git/*" -not -path "./init.sh"  -exec sed -i "s|bartoszmajsak/template-golang|${name}|g" '{}' \;

read -p "Project name: " project
test -z "$name" && {
  die "You must specify project name"
}
sed -i "s|PROJECT_NAME:=template-golang|PROJECT_NAME:=${project}|g" Makefile

read -p "Executable name: " binary
test -z "$name" && {
  die "You must specify executable name"
}
sed -i "s|PROJECT_NAME:=template-golang|BINARY_NAME:=${binary}|g" Makefile
