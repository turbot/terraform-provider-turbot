#!/bin/bash
go clean
go build -o terraform-provider-turbot
terraform init
