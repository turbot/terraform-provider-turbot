#!/bin/bash
if [ -z "$1" ]
  then
    echo "Must provide a version, e.g. 1.0.0-beta.1"
    exit 1
fi

$PWD/scripts/gofmtcheck.sh
mkdir -vp release
env GOOS=linux GOARCH=386 go build -o terraform-provider-turbot_v$1
zip -r release/terraform-provider-turbot_v$1_linux_386.zip terraform-provider-turbot_v$1

env GOOS=linux GOARCH=amd64 go build -o terraform-provider-turbot_v$1
zip -r release/terraform-provider-turbot_v$1_linux_amd64.zip terraform-provider-turbot_v$1

env GOOS=windows GOARCH=386 go build -o terraform-provider-turbot_v$1
zip -r release/terraform-provider-turbot_v$1_windows_386.zip terraform-provider-turbot_v$1

env GOOS=windows GOARCH=amd64 go build -o terraform-provider-turbot_v$1
zip -r release/terraform-provider-turbot_v$1_windows_amd64.zip terraform-provider-turbot_v$1

env GOOS=darwin GOARCH=amd64 go build -o terraform-provider-turbot_v$1
zip -r release/terraform-provider-turbot_v$1_darwin_amd64.zip terraform-provider-turbot_v$1
