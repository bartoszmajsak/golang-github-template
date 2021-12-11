#!/bin/bash

die () {
    echo >&2 "$@"
    exit 1
}

defaultName=$(git remote -v | cut -d':' -f 2 | cut -d'.' -f 1 | uniq)
read -p "Repo name (${defaultName}): " name
test -z "$name" && name=${defaultName}

defaultProject=${defaultName#*/}
read -p "Project name (${defaultProject}}: " project
test -z "project" && project=${defaultProject}

read -p "Executable name: " binary
test -z "binary" && {
  die "You must specify executable name"
}

find . -type f -not -path "./vendor/*" -not -path "./.git/*" -not -path "./init.sh"  -exec sed -i "s|bartoszmajsak/template-golang|${name}|g" '{}' \;
sed -i "s|PROJECT_NAME:=template-golang|PROJECT_NAME:=${project}|g" Makefile
sed -i "s|BINARY_NAME:=binary|BINARY_NAME:=${binary}|g" Makefile
